---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: python 内置模块使用
---
# python 内置模块使用

##### os 模块

| 函数名            | 使用方法                                                     |
| ----------------- | ------------------------------------------------------------ |
| getcwd()          | 返回当前工作目录                                             |
| chdir(path)       | 改变工作目录                                                 |
| listdir(path='.') | 列举指定目录中的文件名                                       |
| mkdir(path)       | 创建单层目录，如目录已存在抛出异常                           |
| makedirs(path)    | 递归创建多层目录，如该目录已存在抛出异常                     |
| remove(path)      | 删除文件                                                     |
| rmdir(path)       | 删除单层目录，如该目录非空抛出异常                           |
| removedirs(path)  | 递归删除目录，从子目录到父目录逐层尝试删除，遇到目录非空则抛出异常 |
| rename(old,new)   | 将文件 old 重命名为 new    |
|  system(command)	  |  运行系统的 shell 命令 |
|walk(top)| 遍历 top 路径下所有的子目录，返回一个三元组：（路径，[包含目录]，[包含文件]|
|os.curdir | 指代当前目录（'.'） |
| os.pardir | 指代上一级目录 ('.') |
| os.sep | 输出操作系统特定的路径分隔符 |
| os.linesep | 当前平台使用的行终止符 |
| os.name | 指代当前使用的操作系统 |

os.walk(top,topdown=True,onerror=None,followlinks=False) 

--top 是要遍历的目录。 

--topdown 是代表要 True 从上而下遍历 False 从下往上遍历 

--onerror 可以用来表示设置当遍历出现错误的处理函数（该函数接收一个 OSError 得到实例作为参数），设置为空则不作处理 

--followlinks 表示是否要跟随目录下的链接去继续遍历，要注意的是，os.walk 不会记录已经遍历的目录，所以跟随链接遍历有可能一直循环调用下去 os.walk 返回的是一个 3 个元素的元组（root,dirs,files）,分别表示遍历的路径名，该路径下的目录列表和该路径下文件列表。注意目录列表和文件列表不是具体路径，需要具体路径（从 root 开始的路径）的话可以用 os.path.join(root,dir) 和 os.path.join(root,dir) 

##### os.path 模块中关于路径常用的函数

| 函数名                      | 使用方法                                                     |
| --------------------------- | ------------------------------------------------------------ |
| basename(path)              | 去掉目录路径，单独返回文件名                                 |
| dirname(path)               | 去掉文件名，单独返回目录路径                                 |
| join(path1,[path2],[ ... ]) | 将 path1,path2 各部分组合成一个路径名                         |
| split(path)                 | 分割文件名与路径，返回 (f_path,f_name) 元组。如果使用完全目录，会将最后一个目录为文件名分离，且不会判断文件或者目录是否存在 |
| splitext(path)              | 分离文件名与扩展名，返回 (f_name,f_extension) 元组             |
| getsize(file)               | 返回指定文件的尺寸，单位是字节                               |
| getatime(file)              | 返回指定文件最近的访问时间（浮点型秒数，可用 time 模块的 gmtime() 或 localtime() 函数换算 |
| getctime(file)              | 返回指定文件的创建时间（浮点型秒数，可用 time 模块的 gmtime() 或 localtime() 函数换算） |
| getmtime(file)              | 返回指定文件最新的修改时间（浮点型秒数，可用 time 模块的 gmtime() 或 localtime() 函数换算） |
| exists(path)                | 判断指定路径（目录或文件）是否存在                           |
| isabs(path)                 | 判断指定路径是否为绝对路径                                   |
| isdir(path)                 | 判断指定路径是否存在且是一个目录                             |
| isfile(path)                | 判断指定路径是否存在且是一个文件                             |
| islink(path)                | 判断指定路径是否存在且是一个符号链接 |
| ismount(path)               | 判断指定路径是否存在且是一个挂载点 |
| samefile(path1,path2) | 判断 path1 和 path2 两个路径是否指向同一个文件 |

##### shutil 模块 文件、文件夹的移动、复制

| 函数名                       | 使用方法                                     |
| ---------------------------- | -------------------------------------------- |
| copy(path1\file,path2\\file) | 从 path1 复制到 path2，也可以复制并重命名新文件 |
| copytree(path1,path2)        | 复制整个目录                                 |
| shutil.rmtree(path)          | 删除文件夹及内容，不管当前目录是否为空       |
| shutil.move(path[file],path2[file]) | 移动文件夹或文件,也可以移动并重命名 |
