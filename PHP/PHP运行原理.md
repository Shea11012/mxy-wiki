## PHP 运行原理

### PHP 的设计理念及特点

PHP 是多进程模型，不同请求间互不干涉，这样保证了一个请求挂掉不会对全盘服务造成影响。PHP 现在也支持多线程模型。PHP 是一门弱类型语言，运行时才会确定并可能发生隐式或显示的类型转换。引擎（Zend）+ 组件（ext）的模式降低内部耦合。中间层（sapi）隔绝 web server 和 PHP。

PHP 四层体系

php-core

![](https://camo.githubusercontent.com/ff6579951d19f3a0282842e8c5ec29436ec30118/68747470733a2f2f7777772e617761696d61692e636f6d2f77702d636f6e74656e742f75706c6f6164732f323031362f30322f7068702d636f72652d343736783530302e706e67)

PHP 从上到下 4 层体系

1. Zend 引擎：

Zend 整体用纯 C 实现，是 PHP 的内核部分，它将 PHP 代码翻译（词法、语法解析等一系列编译过程）为可执行 opcode 处理，并实现相应的处理方法，实现了基本的数据结构（如 hashtable，oo）、内存分配及管理、提供了相应的 api 方法供外部调用，是一切的核心，所有的外围功能均围绕 Zend 实现。

2. Extensions：

围绕着 Zend 引擎，extensions 通过组件式的方式提供各种基础服务，常见的各种内置函数（如 array 系列）、标准库等都是通过 extensions 来实现，用户也可以根据需要实现自己的 extension 已达到功能扩展、性能优化等目的（如贴吧正在使用的 PHP 中间层、富文本解析就是 extension 的典型应用）。

3. Sapi：

Sapi（Server Application Programming Interface）服务端应用编程接口，Sapi 通过一系列钩子函数，使得 PHP 可以和外围交互数据，这是 PHP 非常优雅和成功的一个设计，通过 Sapi 成功将 PHP 本身和上层应用解耦隔离，PHP 可以不再考虑如何针对不同应用进行兼容，而应用本身也可以针对自己的特点实现不同的处理方式。

4. Application：

这层就是平时编写的 PHP 程序，通过不同的 Sapi 方式得到各种各样的应用模式，如通过 webserver 实现 web 应用、在命令行下以脚本方式运行。

#### SAPI

SAPI 通过一系列的接口，使得外部应用可以和 PHP 交换数据，并可以根据不同应用特点实现特定的处理方法，常见的一些 sapi 有：

apache2handler：这是以 Apache 作为 webserver ，采用 mod_php 模式运行时候的处理方式。

cgi：webserver 和 PHP 直接的另一种交互方式，也就是 fastcgi 协议，也是异步 webserver 唯一支持的方式。

cli：命令行调用的应用模式。

#### PHP 的执行流程 opcode

![](https://camo.githubusercontent.com/2765edbcd38c9bb59e0ac5c8217a386637bf0532/68747470733a2f2f7777772e617761696d61692e636f6d2f77702d636f6e74656e742f75706c6f6164732f323031362f30322f323031315f30395f32305f30322e6a7067)

从图上可以看到，PHP 实现了一个典型的动态语言执行过程：拿到一段代码后，经过词法解析、语法解析等阶段后，源程序会被翻译一个个指令（opcodes），然后 ZEND 虚拟机顺次执行这些指令完成操作。PHP 本身是用 C 实现的，因此最终调用的也都是 C 的函数，实际上，可以把 PHP 看做是一个 C 开发的软件。

Opcode 是 PHP 程序执行的最基本单位。一个 opcode 由两个参数（op1，op2）、返回值和处理函数组成。PHP 程序最终被翻译为一组 opcode 处理函数的顺序执行。

常见的几个处理函数：

ZEND_ASSIGN_SPEC_CV_CV_HANDLER：变量分配（\$a = \$b）

ZEND_DO_FCALL_BY_NAME_SPEC_HANDLER：函数调用

ZEND_CONCAT_SPEC_CV_CV_HANDLER：字符串拼接 \$a.$b

ZEND_ADD_SPEC_CV_CONST_HANDLER：加法运算 $a + 2

ZEND_IS_EQUAL_SPEC_CV_CONST：判断相等 \$a ===1

ZEND_IS_IDENTICAL_SPEC_CV_CONST：判断相等 $a === 1

#### HashTable —— 核心数据结构

HashTable 是 Zend 的核心数据结构，在 PHP 里面几乎并用来实现所有常见功能，我们知道 PHP 的数组即是其典型应用，此外，在 zend 内部，如函数符号表、全局变量等也都是基于 hash table 来实现。

PHP 的 hash table 具有如下特点：

支持典型的 key-> value 查询可以当做数组使用添加、删除节点是 O(1) 复杂度 key 支持混合类型：同时存在关联数组索引数组 Value 支持混合类型：`array('string','2332')` 支持线性遍历：如 foreach Zend hash table 实现了典型的 hash 表散列结构，同时通过附加一个双向链表，提供了正向、反向遍历数组的功能。结构如下图：

![](https://www.awaimai.com/wp-content/uploads/2016/02/2011_09_20_03.jpg)

**散列结构：**

Zend 的散列结构是典型的 hash 表模型，通过链表的方式来解决冲突。需要注意的是 Zend 的 hash table 是一个自增长的数据结构，当 hash 表数目满了之后，其本身会动态以 2 倍的方式扩容并重新元素位置。初始大小均为 8。另外，在进行 key->value 快速查找时候，zend 本身还做了一些优化，通过空间换时间的方式加快速度。比如在每个元素中都会用一个变量 `nKeyLength` 标识 key 的长度以作快速判定。

**双向链表：**

Zend hash table 通过一个链表结构，实现了元素的线性遍历，理论上，做遍历使用单向链表就够了，之所以使用双向链表，主要目的是为了快速删除，避免遍历。Zend hash table 是一种复合型的结构，作为数组使用时，即支持常见的关联数组也能够作为顺序索引数字来使用，甚至允许 2 者的混合。

**PHP 关联数组：**

关联数组是典型的 hash table 应用，一次查询过程经过如下几步：

```c
getKeyHashValue h;
index = n & nTableMask;
Bucket *p = arBucket[index];
while(p){
    if((p->h == h) & (p->KeyLength == nKeyLength)) {
        RETURN p->data;
    }
    p = p->next;
}
RETURN FALTURE;
```

**PHP 索引数组：**

索引数组是常见的数组，通过下标访问，例如 `$arr[0]` ，Zend HashTable 内部进行了归一化处理，对于 index 类型 key 同样分配了 hash 值和 nKeyLength。内部成员变量 nNextFreeElement 就是当前分配到的最大 id，每次 push 后自动加一。正是这种归一化处理，PHP 才能实现关联和非关联的混合，由于 push 操作的特殊性，索引 key 在 PHP 数组中先后顺序并不是通过下标大小来决定，而是由 push 的先后决定。`$arr[1] = 2;$arr[2] = 3;`，对于 double 类型的 key，ZendHashTable 会将它当做索引 key 处理。

#### PHP 变量

PHP 是一门弱类型语言，本身不严格区分变量的类型。PHP 在变量申明的时候不需要指定类型。PHP 变量可以分为简单类型（int、string、bool）、集合类型（array、resource、object）和常量（const）。以上所有的变量在底层都是同一种 `zval`

Zval 是 Zend 中的一个非常重要的数据结构，用来标识并实现 PHP 变量，其数据结构如下：

![](https://www.awaimai.com/wp-content/uploads/2016/02/2011_09_20_04.jpg)

#### 引用计数

引用计数在内存回收、字符串操作等地方使用非常广泛。PHP 变量就是引用计数的应用。Zval 的引用计数通过成员变量 `is_ref` 和 `ref_count` 实现，通过引用计数，多个变量可以共享同一份数据，避免频繁拷贝带来的大量消耗。

在进行赋值操作时，zend 将变量指向相同的 `zval` 同时 `ref_count++` ，在 `unset` 操作时，对应的 `ref_count-1` ，只有 `ref_count = 0` 时才会真正执行销毁操作，如果是引用赋值，zend 会修改 `is_ref=1`

 PHP 当试图写入一个变量时，Zend 若发现该变量指向的 `zval` 被多个变量共享，则为其复制一份 `ref_count` 为 `1` 的 `zval` ，并递减原 zval 的 `ref_count` ，这个过程称为 zval 分离。可见，只有在写操作发生时 zend 才进行拷贝操作，因此也叫 `copy-on-write` （写时复制）。

##### 整数和浮点数

整数、浮点数是 PHP 中的基础类型之一，也是一个简单的变量。对于整数和浮点数，在 `zvalue` 中直接存储对应的值，其类型分别是 `long` 和 `double` 。

从 `zvalue` 结构中看出，对于整数类型，和 C 等强类型语言不同，PHP 是不区分 `int、unsigned int、long、 long long` 等型，对于整数只有一种类型就是 `long` 。由此可以看出，在 PHP 里面，整数的取值范围是由编译器位数来决定而不是固定不变的。

对于浮点数，类似整数，也不区分 `float` 和 `double` 而是统一只有 `double` 类型。

在 PHP 中如果整数范围越界，这种情况下会自动转换为 `double` 类型，很多 trick 都是由此产生。

##### 字符和字符串

在 PHP 中，字符串是由指向实际数据的指针和长度结构体组成，和 C++ 的 string 类似。由于通过一个实际变量表示长度，和 C 不同，它的字符串可以是 二进制数据，同时在 PHP 中，求字符串长度 `strlen` 是 0(1) 复杂度。

##### 数组

PHP 的数组通过 Zend HashTable 实现

对一个数组的 `foreach` 是通过遍历 HashTable 中的双向链表完成。对于索引数组，通过 `foreach` 遍历效率比 for 高，省去了 key->value 查找。

##### 资源

资源类型变量是 PHP 中最复杂的一种变量，也是一种复合型结构。

PHP 的 `zval` 可以表示广泛的数据类型，但是对于自定义的数据类型却很难充分描述。由于没有有效的方式描绘这些复合结构，因此没有办法对它们使用传统的操作符。解决这个问题，需要通过一个本质上任意的标识符（label）引用指针，这种方式被称为资源。

`zval` ，对于 resource，lval 作为指针来使用，直接指向资源所在的地址。Resource 可以是任意的复合结构，如mysqli，fsock，memcached等都是资源。

如何使用资源：

- 注册：对于一个自定义的数据类型，要想将它作为资源。首先需要进行注册，zend 会为它分配全局唯一标识。
- 获取一个资源变量：对于资源，zend 维护一个 id->实际数据的 hash_table。对于一个 resource，在 zval 中只记录了它的 id。fetch 的时候通过 id 在 hash_table 中找到具体的值返回。
- 资源销毁：资源的数据类型是多种多样。Zend 本身没有办法销毁它。因此需要用户在注册资源的时候提供销毁函数。当 unset 资源时，zend 调用相应的函数完成析构，同时从全局资源删除它。

资源可以长期驻留，不只是在所有引用它的变量超出作用域之后，甚至是在一个请求结束了并且新的请求产生之后。这些资源称为持久资源，因为它们贯通 SAPI 的整个生命周期持续存在，除非特意销毁。很多情况下，持久化资源可以在一定程度上提高性能。比如常见的 `mysql_pconnect` ，持久化资源通过 `pemalloc` 分配内存，这样请求结束的时候不会释放。

##### 变量作用域

PHP 中的局部变量和全局变量的实现。对于一个请求，任意时刻 PHP 都可以看到两个符号表 symbol_table 和 active_symbol_table ，前者用来维护全局变量，后者是一个指针，指向当前活动的变量符号表。当程序进入到某个函数时，zend 就会为它分配一个符号表 x 同时将 active_symbol_table 指向 a。通过这样的方式实现全局、局部变量的区分。

获取变量值：PHP 符号表是通过 hash_table 实现，对于每个变量都分配唯一标识，获取的时候根据标识从表中找到相应 zval 返回。

函数中使用全局变量：在函数中，我们可以通过显式申明 global 来使用全局变量。在 `active_symbol_table` 中创建 `symbol_table` 中同名变量的引用，如果 `symbol_table` 中没有同名变量则会创建。