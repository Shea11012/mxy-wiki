---
title: "PHP 短路"
created: 2021-12-03
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - php
confidence: INFERRED
---
# PHP 短路

只有逻辑与和逻辑或存在短路

```php
$a = 1;
$a || $b = 100;	//$b = 100
$a && $b = 100;	//$b = 100,如果$a = 0,则$b
```



