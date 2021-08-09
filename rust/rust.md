---
tags: ["rust"]
---
# rust

rust 中变量默认不可变，变量需要可变关键字 **mu**
```rust
let foo = "a";	// 默认不可变
let mu bar = "b"; // 变量可变
```
let 与 mu 的区别：
- 可以使用 let 关键字重复声明一个变量，每次 let 声明都是创建了一个新变量，可以改变值类型，但复用了这个名字

```rust
let a = "";
let a = 1;
let a = 1.2
// 以上使用使用相同变量隐藏了一个变量，行为看似是可以对一个变量赋值不同的类型

// mut 关键字
let mut a = " ";
a = 2;	// 会报错，a 变量可变，但是不能改变其类型
```

## 表达式和语句
**rust 中表达式有返回值，语句没有返回值**
```rust
let y = {
	let x = 3;
	x + 1	// 因为没有分号，所以是表达式会有返回值
}
println!("the value of y is: {}",y);
```

## 枚举
```rust
enum IpAddr {
	V4,
	V6,
}

enum IpAddr {
	V4(String),
	V6(String),
}

enum IpAddr {
	V4(u8,u8,u8,u8),
	V6(String),
}

enum Message {
	Quit,
	Move {x: i32, y: i32},
	Write(String),
	ChangeColor(i32,i32,i32),
}

// 在枚举上定义方法
impl Message {
	fn call(&self) {
		// do something
	}
}
```

**rust 没有 null，使用 option 替代**
```rust
enum Option<T> {
	Some(T),
	None,
}
```

## match 控制流运算符
```rust
let value = 8;
match value {
	1 => println!("one"),
	2 => println!("two"),
	3 => println!("three"),
	_ => (),
}
```

### if let
```rust
let value = 8;
match value {
	8 => println!("eight"),
	_ => (),
}


// 以上可以改为
if let 8 = value {
	println!("eight");
}
```

## 所有权
> 使得 rust 无需垃圾回收也可保障内存安全

### 引用与借用
[[../计算机基础/栈和堆#堆栈|堆栈简介]]
#### 所有权规则
1. rust 中的每一个值都有一个被称为所有者（owner）的变量
2. 值在任一时刻有且只有一个所有者
3. 当所有者离开作用域，这个值将被丢弃

#### 引用与借用
```rust
fn main() {
    let s1 = String::from("hello");

    let len = calculate_length(&s1);

    println!("The length of '{}' is {}.", s1, len);
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```

`&` 是引用符号，允许使用值但不获取其所有权，rust 将获取引用作为函数参数称为借用（borrowing）。
rust 默认不允许修改借用的值

#### 可变引用
```rust
fn main() {
    let mut s = String::from("hello");

    change(&mut s);
}

fn change(some_string: &mut String) {
    some_string.push_str(", world");
}
```

可变引用有一个限制：在特定作用域中的特定数据只能有一个可变引用
```rust
let mut s = String::from("hello");

// 以下代码会报错
let r1 = &mut s;
let r2 = &mut s;

println!("{}, {}", r1, r2);
```
这种限制的优点是 rust 可以在编译时就避免数据竞争。数据竞争的竟态条件，可由：
- 两个或更多指针同时访问同一数据
- 至少有一个指针被用来写入数据
- 没有同步数据访问机制

```rust

#![allow(unused)]
fn main() {
let mut s = String::from("hello");

{
    let r1 = &mut s;

} // r1 在这里离开了作用域，所以我们完全可以创建一个新的引用

let r2 = &mut s;
}

```

同时也不能在拥有不可变引用的同时拥有可变引用，多个不可变引用是允许同时存在的。
一个引用的作用域从声明的地方开始一直持续到最后一次使用为止。
```rust
let mut s = String::from("hello");

let r1 = &s; // 没问题
let r2 = &s; // 没问题
println!("{} and {}", r1, r2);
// 此位置之后 r1 和 r2 不再使用

let r3 = &mut s; // 没问题
println!("{}", r3);

```

## 结构体
定义结构体
```rust
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}
```

默认创建的结构体是不可变的，如果需要更改结构体内的某个值，需要使结构体可变

```rust
let mut user1 = User {
	email: String::from("xx.com"),
	username: String::from("sfsf"),
	active: true,
	sign_in_count: 1,
};
```

创建结构体简写语法
```rust

struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

fn build_user(email: String, username: String) -> User {
    User {
        email,
        username,
        active: true,
        sign_in_count: 1,
    }
}

```

结构体更新语法
```rust

#![allow(unused)]
fn main() {
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

let user1 = User {
    email: String::from("someone@example.com"),
    username: String::from("someusername123"),
    active: true,
    sign_in_count: 1,
};

let user2 = User {
    email: String::from("another@example.com"),
    username: String::from("anotherusername567"),
    ..user1
};
}

```


## 包管理、模块
1. 将模块放在 src 目录内 

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210629100115.png)

2. 创建一个本地库，导入本地库

![[rust/Pasted image 20210629111127.png]]

首先在 src 的同级目录下创建了一个 utils 目录，作为一个外部仓库导入到 src 中

1. 首先修改根目录的 cargo.toml 文件

```toml
utils = {path="utils",version="0.1.0"}
```

2. 在 utils 内创建好包后，需要在 lib.rs 中指定导出

![[rust/Pasted image 20210629111511.png]]

对于是在 utils 内的子包，需要在子包内建一个 mod.rs 指定需要导出的模块


## 泛型、trait、生命周期

### 泛型
```rust
struct Point<T> {
	x: T,
	y: T,
}

impl<T> Point<T> {
	fn mixup<V,W>(&self,other: Point<V,W>) -> Point<T,W> {
		Point{
			x: self.x,
			y: other.y
		}
	}
}
```

### trait
> trait 类似其他语言的 interface

定义一个 trait
```rust
pub trait Summary {
	fn summarize(&self) -> String;
}
```

为一个结构体实现 trait
```rust
use std::fmt::Display;

struct Point {
	x: i32,
	y: i32,
}

impl Display for Point {
	fn cmp_display(&self) {
		// do something
	}
}
```

为一个泛型实现 trait
```rust
use std::fmt::Display;

struct Point<T> {
	x: T,
	y: T,
}

impl<T: Display> Point<T> {}	// 只实现一种 trait

impl<T: Display + Partialord> Point<T> {}

struct Pair<V,W> {}

impl<V,W> Point<V,W>
	where V: Display + PartialOrd,	// 这是一种实现多个 trait 的简便写法
		  W: some_trait + some_trait
{
	fn cmp_display(&self) {
		// do something
	}
}

```

##  生命周期
**生命周期参数名称必须以 `'` 开头**
```rust
&i32	// 引用
&'a i32	// 带有显示生命周期引用
&'a mut i32	// 带有显示生命周期的可变引用
```

函数生命周期
> 这个例子确保函数输入参数和输出参数生命周期一致
```rust
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}
```

结构体生命周期
```rust
#[derive(Debug)]
struct Earth {
    location: String,
}

#[derive(Debug)]
struct Dinosaur<'a> {
    location: &'a Earth,
    name: String,
}

fn main() {
    let new_york = Earth{
        location: String::from("New York, NY"),
    };

    let t_rex = Dinosaur{
        location: &new_york,
        name: "T Rex".to_string(),
    };

    println!("{:?}",t_rex);
}
```
> 当一个结构体涉及到引用时，需要给它添加生命周期标识，会使得编译器理解 dinosaur 的生命周期不能比 earth 长。

## 闭包
> 闭包以一对 | 开始，在 |params| 中指定参数

对于闭包不用定义类型和返回类型，编译器会自动推断，也可以按照函数的方式定义类型

```rust
let add = |a,b| {
	a + b
}
```

使用结构体存放闭包时，则需要定义闭包的类型
```rust
struct Cacher<T>
	where T: Fn(u32) -> u32
{
	calculation: T,
	value: Option<u32>,
}
```

闭包从环境中捕获一个值，会产生内存开销。rust 提供了三种的闭包方式：
- FnOnce 捕获周围的变量，同时闭包会获取其所有权并在定义闭包时将其移动进闭包。因为获取了所有权则表示该闭包只能被调用一次。
- FnMut 获取可变的借用值所以可以改变环境
- Fn 从周围环境中获取不可变的借用值，相当于完整拷贝了一份

当创建闭包时，rust 会根据使用情况推断希望闭包如何引用环境。
1. 当闭包被调用至少一次，所以所有闭包都实现了 FnOnce
2. 捕获到的变量并没移动它们的所有权到闭包内时就实现了 FnMut
3. 不需要对被捕获的变量进行可变访问时就实现了 Fn

**如果希望闭包强制获取其使用的环境所有权，关键字 move**
```rust
let add = move |x,y| x + y;
```

## 智能指针
一个类型实现了 Deref trait 则允许重载 \* 解引用运算符
`Box<T>` ：允许将一个值放在堆上，使用场景：
- 当有一个在编译时未知大小的类型
- 当有大量数据并希望在确保数据不被拷贝的情况下转移所有权
- 当希望拥有一个值只关心它的类型是否实现了特定 trait 而不是其具体类型时