---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: PHP实用的示例
---
计算还差多少天过生日

```php
$birth = '5-28';
$year = date('Y-');
$datetime1 = new DateTime(date('Y-m-d'));	//获取当前时间
$datetime2 = new DateTime($year.$birth);	//将当前的年份和生日拼接起来
$interval = $datetime1->diff($datetime2);	//当前时间和生日作比较
$day = $interval->format('%R%a');
switch ($day){
    case $day<0:
        echo '今年生日已经过了，⊙﹏⊙∥';
        break;
    case $day==0:
        echo '生日就在今天，*^____^*';
        break;
    default:
        echo '距离生日还差',abs($day),'天';
}
```

