使用工具 <font color="red">fio</font>
fio 参数
- filename：测试文件名称，文件需要在待测试盘
- direct=1：测试过程绕过机器自带buffer
- rw=randwrite：测试随机写
- rw=randrw：测试随机写和读
- bs=16K：单次IO的块文件大小
- size=300MB：测试文件大小
- numjobs=10：测试线程
- runtime=60：测试时间
- ioengine=psync：指定io引擎
- rwmixwrite=30：在混合读写模式下，写占30%
- group_reporting：显示结果

随机读
```shell
fio -filename=/tmp/test_randread -direct=1 -iodepth 1 -thread -rw=randread -ioengine=psync -bs=16k -size=500MB -numjobs=10 -runtime=60 -group_reporting -name=mytest
```

随机写
```shell
fio -filename=/tmp/test_randwrite -direct=1 -iodepth 1 -thread -rw=randwrite -ioengine=psync -bs=16k -size=500MB -numjobs=10 -runtime=60 -group_reporting -name=mytest
```

顺序写
```shell
fio -filename=/temp/test_write -direct=1 -iodepth 1 -thread -rw=write -ioengine=psync -bs=16k -size=300MB -numjobs=10 -runtime=60 -group_reporting -name=mytest
```

顺序读
```shell
fio -filename=/temp/test_read -direct=1 -iodepth 1 -thread -rw=read -ioengine=psync -bs=16k -size=300MB -numjobs=10 -runtime=60 -group_reporting -name=mytest
```

混合随机读写
```shell
fio -filename=/temp/test_randrw -direct=1 -iodepth 1 -thread -rw=randrw -rwmixread=70 -ioengine=psync -bs=16k -size=300MB -numjobs=10 -runtime=60 -group_reporting -name=mytest
```