---
created: 2026-06-02
updated: 2026-06-02
status: stable
type: practice
tags: [mysql/lock, mysql/innodb, golang, gorm, deadlock]
related:
  - "锁"
---

# RC 级别死锁避免 —— Go/GORM 实践

> 配合 [[锁#5.7 RC 级别下依然会遇到的死锁|锁 §5.7]] 的理论背景阅读。
> RC 禁用了 Gap Lock，死锁概率大幅下降，但不会降为零。
> 本文用 GORM 代码演示 6 个实战策略。

---

## 策略一：统一加锁顺序 —— ORDER BY id

RC 下虽然没有 Gap Lock，但两个事务**反向扫描**同一批行仍会死锁。
核心原则：所有批量操作统一按主键升序加锁。

**❌ 错误：无 ORDER BY，扫描顺序由优化器决定**

```go
func badBatchUpdate(db *gorm.DB, status int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 无序扫描，两个并发事务可能一正一反 → 死锁
        return tx.Model(&Order{}).
            Where("status = ?", 0).
            Updates(map[string]interface{}{"status": status}).Error
    })
}
```

**✅ 正确：按 id 升序 + 分批处理**

```go
func goodBatchUpdate(db *gorm.DB, status int, batchSize int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        var ids []int64
        // 第一步：按 id 升序锁定待处理行
        if err := tx.Model(&Order{}).
            Where("status = ?", 0).
            Order("id ASC"). // ← 关键：统一加锁顺序
            Limit(batchSize).
            Pluck("id", &ids).Error; err != nil {
            return err
        }
        if len(ids) == 0 {
            return nil
        }
        // 第二步：批量更新已锁定的行
        return tx.Model(&Order{}).
            Where("id IN ?", ids).
            Update("status", status).Error
    })
}
```

---

## 策略二：添加索引，避免全表扫描锁

RC 下最常见的死锁根因：WHERE 条件无索引 → 全表扫描 → 锁住大量行。

```go
// ❌ 反例：status 无索引，两个事务全表扫描，方向相反 → 死锁
// T1: UPDATE orders SET status=1 WHERE user_id=100;
// T2: UPDATE orders SET status=1 WHERE user_id=200;

// ✅ 正解：user_id 有组合索引 (user_id, status)，只锁匹配行
func updateUserOrders(db *gorm.DB, userID int64, newStatus int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        return tx.Model(&Order{}).
            Where("user_id = ? AND status = ?", userID, 0).
            Update("status", newStatus).Error
    })
}
```

> [!important] 索引是第一道防线
> 在 RC 下，无索引的 UPDATE/DELETE 会锁**全表所有行**（因为逐行扫描判断是否匹配）。
> 即使两个事务更新的行完全不重叠，扫描方向相反也会死锁。
> 解法只有一个：**给 WHERE 条件加索引**。

---

## 策略三：缩短事务 —— RPC 移到事务外

持锁时间越长，与其他事务冲突的窗口越大。

```go
// ❌ 错误：事务内做 RPC，持锁时间拉长
func badCreateAndNotify(db *gorm.DB, order *Order) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        // RPC 在事务内 → 其他事务等锁时间被拉长
        return callExternalService(order)
    })
}

// ✅ 正确：先提交事务，再做 RPC
func goodCreateAndNotify(db *gorm.DB, order *Order) error {
    // 1. 事务内只做 DB 操作
    if err := db.Transaction(func(tx *gorm.DB) error {
        return tx.Create(order).Error
    }); err != nil {
        return err
    }
    // 2. 事务提交后，再做外部调用
    return callExternalService(order)
}
```

---

## 策略四：NOWAIT / SKIP LOCKED（MySQL 8.0+）

### NOWAIT：遇到锁立即报错，不等待

```go
func lockOrderNowait(db *gorm.DB, orderID int64) (*Order, error) {
    var order Order
    err := db.Transaction(func(tx *gorm.DB) error {
        return tx.Raw(
            "SELECT * FROM orders WHERE id = ? FOR UPDATE NOWAIT", orderID,
        ).Scan(&order).Error
    })
    if err != nil {
        if strings.Contains(err.Error(), "3572") {
            // MySQL Error 3572: lock could not be acquired (NOWAIT)
            return nil, fmt.Errorf("order %d is locked, retry later", orderID)
        }
        return nil, err
    }
    return &order, nil
}
```

### SKIP LOCKED：跳过已被锁的行（队列消费场景）

```go
func claimPendingOrders(db *gorm.DB, limit int) ([]Order, error) {
    var orders []Order
    err := db.Transaction(func(tx *gorm.DB) error {
        return tx.Raw(`
            SELECT * FROM orders
            WHERE status = 0
            ORDER BY id ASC
            LIMIT ?
            FOR UPDATE SKIP LOCKED
        `, limit).Scan(&orders).Error
    })
    return orders, err
}
```

> [!tip] 天然避免死锁
> `SKIP LOCKED` 跳过冲突行不等待，两个事务永远不会互相阻塞，从根本上消除了死锁可能。
> 适用于：任务队列、消息消费、批量分配等场景。

---

## 策略五：死锁重试（兜底）

MySQL 检测到死锁时自动回滚代价较小的事务，返回 **Error 1213**。
应用层捕获后随机退避重试：

```go
func withDeadlockRetry(fn func() error, maxRetries int) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        lastErr = fn()
        if lastErr == nil {
            return nil
        }
        // MySQL Error 1213: Deadlock found when trying to get lock
        if strings.Contains(lastErr.Error(), "Deadlock") ||
            strings.Contains(lastErr.Error(), "1213") {
            backoff := time.Duration(50+rand.Intn(100*(i+1))) * time.Millisecond
            time.Sleep(backoff)
            continue
        }
        return lastErr // 非死锁错误，直接返回
    }
    return fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}
```

**使用示例：**

```go
func processOrderWithRetry(db *gorm.DB, orderID int64) error {
    return withDeadlockRetry(func() error {
        return db.Transaction(func(tx *gorm.DB) error {
            var order Order
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
                Where("id = ?", orderID).
                Order("id ASC"). // 统一排序
                First(&order).Error; err != nil {
                return err
            }
            order.Status = 1
            return tx.Save(&order).Error
        })
    }, 3)
}
```

---

## 策略六：转账场景 —— 排序加锁

两个账户互相转账，如果各锁各的就死锁。解法：**永远按 id 从小到大加锁**。

```go
func transfer(db *gorm.DB, fromID, toID int64, amount int64) error {
    // 确保加锁顺序一致：小 ID 先锁
    first, second := fromID, toID
    if first > second {
        first, second = second, first
    }

    return withDeadlockRetry(func() error {
        return db.Transaction(func(tx *gorm.DB) error {
            var accounts []Order
            // 按 id 升序锁定两个账户 ← 关键
            if err := tx.Where("id IN ?", []int64{first, second}).
                Order("id ASC").
                Clauses(clause.Locking{Strength: "UPDATE"}).
                Find(&accounts).Error; err != nil {
                return err
            }
            if len(accounts) != 2 {
                return fmt.Errorf("account not found")
            }

            // 根据实际 fromID/toID 扣减/增加
            for i := range accounts {
                if accounts[i].ID == fromID {
                    if accounts[i].Amount < amount {
                        return fmt.Errorf("insufficient balance")
                    }
                    accounts[i].Amount -= amount
                } else {
                    accounts[i].Amount += amount
                }
            }
            return tx.Save(&accounts).Error
        })
    }, 3)
}
```

> [!important] 排序是消除死锁的万能钥匙
> 不管多少个事务、操作多少行，只要保证**加锁顺序一致**（按主键升序），就不会出现循环等待。
> 这是唯一一个能从理论上完全消除死锁的策略。

---

## 策略速查

| # | 策略 | 适用场景 | 能否完全消除死锁 |
|---|------|---------|:---:|
| 1 | ORDER BY id 统一顺序 | 批量 UPDATE/DELETE | ✅ |
| 2 | 加索引 | WHERE 条件列 | 减少概率 |
| 3 | 缩短事务 | 含 RPC/计算的事务 | 减少概率 |
| 4 | NOWAIT / SKIP LOCKED | 队列消费、热点行竞争 | ✅ |
| 5 | 死锁重试 | 兜底，所有事务都应套一层 | 容忍并恢复 |
| 6 | 排序加锁 | 转账、双向操作 | ✅ |

---

## 连接配置

```go
func initDB() *gorm.DB {
    dsn := "user:pass@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
    db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    // 设置会话级隔离级别为 RC（也可在 my.cnf 全局设置）
    db.Exec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
    return db
}
```

---

## 相关页面

- [[锁]] — InnoDB 锁机制完整笔记
- [[事务]] — 隔离级别、幻读、MVCC
