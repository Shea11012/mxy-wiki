---
tags: ["memory ordering", "atomic"]
date created: 2022-12-02 10:06
date modified: 2022-12-06 18:15
---
## 原子操作
```rust
use std::sync::atomic::{
    AtomicU32,
    Ordering,
};

fn main() {
    let x = AtomicU32::new(0);
    x.fetch_add(1, Ordering::Relaxed);
    x.fetch_add(1, Ordering::Relaxed);
    println!("{}",x.load(Ordering::Relaxed));
}
```


## 内存顺序

内存顺序是 CPU 在访问内存时的顺序，会受到一下影响：

- 代码中的先后顺序
- 编译器优化导致在编译阶段发生改变
- 运行阶段因 CPU 的缓存机制导致顺序被打乱

rust 提供了 5 种内存排序

```rust
pub enum Ordering {
	Relaxed,
	Release,
	Acquire,
	AcqRel,
	SeqCst,
}
```

### Relaxed

Relaxed 不施加任何限制，CPU 和编译器可以自由 reorder，仅保证原子性

```rust
let x: &'static _ =  Box::leak(Box::new(AtomicUsize::new(0)));
let y: &'static _ =  Box::leak(Box::new(AtomicUsize::new(0)));

let t1 = spawn(move || {
	let r1 = y.load(Ordering::Relaxed);
	x.store(r1,Ordering::Relaxed);
	r1
});

let t2 = spawn(move ||{
	let r2 = x.load(Ordering::Relaxed);
	y.store(Ordering::Relaxed);
	r2
});

let r1 = t1.join().unwarp();
let r2 = t2.join().unwrap();
```

上面的代码，r1 可能等于 r2，因为 Relaxed 不保证 happen before，就导致 CPU 可能进行指令重拍

### Release & Acquire

Release： 针对写操作，任何读写操作都不能重排到该写操作之后，并且所有当前线程中在该原子操作之前的写操作都对另一个使用 Acquire 读操作的线程可见

Acuqire：针对读操作，任何读写操作都不能被重排到该读操作之前，并且当前线程可以看到另一个线程对同一个变量使用 Release 写操作之前的所有写操作。

```rust
// main
let data: &'static _ = Box::leak(Box::new(AtomicUsize(0)));
let flag: &'static _ = Box::leak(Box::new(AtomicBool(false)));

// thread 1
let t1 = spawn(move || {
	data.store(1,Ordering::Relaxed); // A
	flag.store(true, Ordering::Release); // B
});

// thread 2
let t2 = spawn(move || {
	while !flag.load(Ordering::Acquire) {} // C
	assert_eq!(data.load(Ordering::Relaxed,1); // D
});
```

thread1，根据 Release 的规则，A 不会重排到 B 后面

thread2，根据 Acquire 规则，D 不会重排到 C 的前面

如果不使用 Release-Acquire Ordering 就可能导致指令重排

### AcqRel

主要用于 read-modify-write 操作，如果 compare_and_swap，表明它同时具有 Acquire 和 Release 语义，读的部分用 Acquire，写的部分用 Release。

### SeqCst

Sequential Consistent，对读使用 Acquire，对写使用 Release，对 read-modify-write 使用 AcqRel 的基础上再保证所有线程看到所有使用了 SeqCst 的操作是同一个顺序。

链接：
- http://senlinzhan.github.io/2017/12/04/cpp-memory-order/
- https://rustmagazine.github.io/rust_magazine_2022/Q1/contribute/atom.html
- [Crust of Rust: Atomics and Memory Ordering](https://www.youtube.com/watch?v=rMGWeSjctlY) 