Pickle模块中最常用的的函数为：

1. pickl.dump(obj,file,[protocol])：将obj对象序列化存入已经打开的file中

   - obj：想要序列化的obj对象
   - file：文件名称
   - protocol：序列化使用的协议。如果该项省略，则默认为0。如果为负值或HIGHEST_PROTOCOL，则使用最高的协议版本
2. pickle.load(file)：将file中的对象序列化读出
   -  file：文件名称
3. pickle.dumps(obj,[protocol])：将obj对象序列化为string形式，而不是存入文件
   -  obj：想要序列化的obj对象
   -  protocol：序列化使用的协议。如果该项省略，则默认为0。如果为负值或HIGHEST_PROTOCOL，则使用最高的协议版本
4. pickle.loads(string)：从string中读出序列化前的obj对象
   - string：文件名称

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

