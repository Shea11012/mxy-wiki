---
tags: 
date created: 2022-03-26 11:19
date modified: 2023-03-14 04:04
title: GC
---
![[栈和堆#栈]]
![[栈和堆#堆]]

## GC

GC 的主要对象是堆内存

### GC 方法

#### 手动内存管理

语言不会帮你管理内存，需要开发者自己创建和释放内存，如：C、C++

#### 垃圾回收

通过释放无用的内存自动管理堆内存，GC 会导致 STW（stop the world）。多数语言都实现了 GC，如：JVM（Java、Scala、Groovy、Kotlin），Javascript、C#、Golang、OCaml、Ruby
- Mark & Sweep GC：跟踪 GC，分为两步第一步标记存活（即被引用）的对象，第二部释放无用的对象。JVM、Go、Javascript
- Reference Counting GC：每个对象的引用都有一个计数器进行 +1 或 -1，当计数器变为 0 时则释放。PHP、Perl、Python

#### Resource Acquisition is Initialization (RAII)

当一个对象分配内存时绑定了它的生命周期，在生命周期结束后会被释放。如：Rust

#### Automatic Reference Counting(ARC)

与引用计数类似，但 ARC 在编译器插入指令并计数当计数器为 0 时则释放内存，ARC 回收内存不会造成 STW，但是如引用计数一样不能处理循环引用。Clang、Objective C、Swift

#### Ownership

将 RAII 与 Ownership 模型结合，任何值都必须有一个变量作为它的 owner，当 owner 超出了值范围时，无论是在堆或者栈都将会被释放。Rust