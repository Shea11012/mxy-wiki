### 解压可迭代对象赋值给多个变量

```python
record = ('dave','xx@example.com','773-555-1212', '847-555-1212')
name,email,*phone_number = record
# name = 'dave' email = 'xx@example.com' phone_number = ['773-555-1212', '847-555-1212']
```

如果解压一些未曾使用到的元素，可以使用约定好的标记 `*_或 *ign`

### 字典的运算

关键是使用 `zip` 函数将对象中对应的元素打包成一个个元组

```python
prices = {
    'ACME': 45.23,
    'AAPL': 612.78,
    'IBM': 205.55,
    'HPQ': 37.20,
    'FB': 10.75
}

min_price = min(zip(prices.values(),prices.keys()))

values,keys = zip(*zip(prices.value(),prices.keys()))


# 查找两字典中寻找相同点
a = {
    'x' : 1,
    'y' : 2,
    'z' : 3
}

b = {
    'w' : 10,
    'x' : 11,
    'y' : 2
}

a.keys() & b.keys() # {'x','y'}
a.keys() - b.keys() # {'z'}
a.items() & b.items() # {('y',2)}
```

### 序列中出现次数最多的元素

`collections.Counter`

```python
words = [
    'look', 'into', 'my', 'eyes', 'look', 'into', 'my', 'eyes',
    'the', 'eyes', 'the', 'eyes', 'the', 'eyes', 'not', 'around', 'the',
    'eyes', "don't", 'look', 'around', 'the', 'eyes', 'look', 'into',
    'my', 'eyes', "you're", 'under'
]
from collections import Counter
word_counts = Counter(words)
# 出现频率最高的3个单词
top_three = word_counts.most_common(3)
print(top_three)
# Outputs [('eyes', 8), ('the', 5), ('look', 4)]
```

### 通过某个关键字排序一个字典列表

使用 `operator` 模块的 `itemgetter` 函数

```python
rows = [
    {'fname': 'Brian', 'lname': 'Jones', 'uid': 1003},
    {'fname': 'David', 'lname': 'Beazley', 'uid': 1002},
    {'fname': 'John', 'lname': 'Cleese', 'uid': 1001},
    {'fname': 'Big', 'lname': 'Jones', 'uid': 1004}
]

from operator import itemgetter
rows_by_fname = sorted(rows, key=itemgetter('fname'))
rows_by_uid = sorted(rows, key=itemgetter('uid'))
print(rows_by_fname)
print(rows_by_uid)
```

### 通过某个字段将记录分组

`itertools.groupby()`

```python
rows = [
    {'address': '5412 N CLARK', 'date': '07/01/2012'},
    {'address': '5148 N CLARK', 'date': '07/04/2012'},
    {'address': '5800 E 58TH', 'date': '07/02/2012'},
    {'address': '2122 N CLARK', 'date': '07/03/2012'},
    {'address': '5645 N RAVENSWOOD', 'date': '07/02/2012'},
    {'address': '1060 W ADDISON', 'date': '07/02/2012'},
    {'address': '4801 N BROADWAY', 'date': '07/01/2012'},
    {'address': '1039 W GRANVILLE', 'date': '07/04/2012'},
]

from operator import itemgetter
from itertools import groupby

# Sort by the desired field first
rows.sort(key=itemgetter('date'))
# Iterate in groups
for date, items in groupby(rows, key=itemgetter('date')):
    print(date)
    for i in items:
        print(' ', i)
```

### 映射名称到序列元素

有时通过下标访问列表或者元组中元素的代码，会显得难以阅读。

可以使用 `collections.nanedtuple()` 会创建一个类似元组的对象，可以进行类的操作，也可进行元组的操作

```python
from collections import namedtuple
Subscriber = namedtuple('Subscriber',['addr','joined'])
sub = Subscriber('jones@example.com','2320-13-2')
sub.addr # jones@example.com
sub.joined # 2320-13-2
```



