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

## 定义结构和实现结构体方法
```rust
struct User {
	email: String,
	address: String,
	active: true,
}

impl User {
	fn toggleActive(&self) -> bool {
		!self.active
	}
}
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