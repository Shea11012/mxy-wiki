---
title: "数据存储 pickle模块"
created: 2021-12-03
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - python
confidence: INFERRED
---
Pickle 模块中最常用的的函数为：

1. `pickle.dump(obj, file, protocol=None)`：将对象写入二进制文件对象
   - obj：要序列化的对象
   - file：支持写入 bytes 的二进制文件对象
   - protocol：序列化协议；省略时使用 `pickle.DEFAULT_PROTOCOL`，其值随 Python 版本变化。传负值等价于使用 `pickle.HIGHEST_PROTOCOL`
2. `pickle.load(file)`：从二进制文件对象读取并反序列化一个对象
   - file：支持读取 bytes 的二进制文件对象
3. `pickle.dumps(obj, protocol=None)`：将对象序列化为 bytes，而不是写入文件
   - obj：要序列化的对象
   - protocol：省略时使用 `pickle.DEFAULT_PROTOCOL`；传负值等价于使用 `pickle.HIGHEST_PROTOCOL`
4. `pickle.loads(data)`：从包含 pickle 数据的 bytes-like 对象反序列化对象
   - data：包含 pickle 数据的 bytes-like 对象

```python
import pickle
dataList = [
    [1,1,'yes'],
    [1,1,'yes'],
    [1,0,'no'],
    [0,1,'no'],
    [0,1,'no'],
           ]
dataDic = {0:[1,2,3,4],
          1:('a','b'),
          2:{'c':'yes','d':'no'}}
#使用dump()将数据序列化到文件中
fw = open('dataFile.txt','wb')
pickle.dump(dataList,fw,-1)
pickle.dump(dataDic,fw)
fw.close()

#使用load()将数据从文件序列化读出
fr = open('dataFile.txt','rb')
data1 = pickle.load(fr)
print(data1)
data2 = pickle.load(fr)
print(data2)
fr.close()

#使用dumps()和loads()举例
p = pickle.dumps(dataList)
print(pickle.loads(p))
p = pickle.dumps(dataDic)
print(pickle.loads(p))
```

> [!warning]
> `pickle` 反序列化可能执行任意代码。只能加载可信且未被篡改的数据；处理不可信输入时应选择 JSON 等更安全的格式。

## 核验依据

- [Python 官方文档：pickle](https://docs.python.org/3/library/pickle.html)

