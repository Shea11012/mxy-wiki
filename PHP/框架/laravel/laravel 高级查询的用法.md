### where 的闭包用法

```php
// 在闭包内关联其它表
where('id',function($query){
    // select 指明需要查询的它表数据
    // from 指定它表
    $query->select('id')->from('post'); // from 里的表名是数据库的表名，不是模型名
})
```

```php
$builder = Product::query()->where('on_sale',true);
$builder->where(function ($query) use ($like) {
                $query->where('title','like',$like)
                    ->orWhere('description','like',$like)
                    ->orWhereHas('skus',function ($query) use ($like) {
                        $query->where('title','like',$like)
                            ->orWhere('description','like',$like);
                    });
            });
// 上面的where 用法目的是在查询条件的两边加上 (),最终执行 SQL 语句如下
select * from products where on_sale = 1 and ( title like xxx or description like xxx )
// 如果不使用 where 的 闭包函数，如下
    $like = '%'.$search.'%';
$builder->where('title', 'like', $like)
    ->orWhere('description', 'like', $like)
    ->orWhereHas('skus', function ($query) use ($like) {
        $query->where('title', 'like', $like)
            ->orWhere('description', 'like', $like);
    });
则 sql 会变成如下：
    select * from products where on_sale = 1 and title like xxx or description like xxx
```

