---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel 商城开发复习
---
#### 商品类目

| 字段名称     | 描述               | 类型               | 索引 |
| ------------ | ------------------ | ------------------ | ---- |
| id           | 自增 ID             | unsigned int       | 主键 |
| name         | 类目名称           | varchar            | 无   |
| parent_id    | 父类目 ID           | unsigned int，null | 外键 |
| is_directory | 是否拥有子类目     | tinyint            | 无   |
| level        | 当前类目等级       | unsigned int       | 无   |
| path         | 该类目所有父类目 ID | varchar            | 无   |

path 字段意义：

1. 查询一个类目的所有祖先类目，需要递归地去逐级查询父类目，会产生较多的 SQL 查询，影响性能
2. 查询一个类目的所有后代类目，同样需要递归逐级查询父类目，同样会产生很多 SQL 查询，影响性能
3. 判断两个类目是否有祖孙关系，需要从层级低的类目逐级往上查，性能低下

path 字段保存该类目的所有祖先类目 ID，并用 `-` 分隔



laravel 提供了一种 viewComposer 的功能，可以在不修改控制器的情况下直接向指定的模板文件注入变量

viewcomposer 类 通常放在 `app/Http/ViewComposers` 目录下

```php
class CategoryTreeComposer
{
    protected $categoryService;
    
    // 此处 laravel 会自动依赖注入
    public function __construct(CategoryService $categoryService)
    {
        $this->categorySerivce = $categoryService;
    }
    
    // 当渲染指定模板时， laravel 会调用 compose 方法
    public function compose（View $view)
    {
        $view->with('categoryTree',$this->categoryService->getCategoryTree());
    }
}

// app/Providers/AppServiceProvider

public function boot()
{
 // 当laravel 渲染 products.index products.show 时就会使用 CategoryTreeComposer 注入这个类目树的变量
    \View::composer(['products.index','products.show'],CategoryTreeComposer::class);
}
```

