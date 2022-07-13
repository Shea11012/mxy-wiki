---
tags: ["rust"]
date created: 2021-06-25 21:32
date modified: 2022-01-23 15:09
title: rust
---
# rust
[rust](https://rustup.rs/) 官网安装 rust 编译器和工具

预加载类型：不需要导包就可以直接使用的类型,[https://doc.rust-lang.org/std/prelude/](https://doc.rust-lang.org/std/prelude/) ^f4bc7d
使用 `unimplemented!` 宏，避免报错
```rust
fn dada() -> Result<bool,FileError> {
	unimplemented!();
}
```

## rust 工具
- 更新 rust 工具：`rustup update`
- 更新 rustup：`rustup self update`
- 展示当前工具链：`rustup show`
- 切换 rust 版本：`rustup override set nightly`
- 安装指定版本：`rustup install nightly-2016-06-03`

## rust 概要
- bool
- char
- isize：有符号类型；如：i8、i16、i32、i64、i128
- usize：无符号类型；如：u8、u16、u32、u64、u128
- f32
- f64
- `[T;N]`：固定大小的数组；元素类型为 T，N 在编译时确定
- `[T]`：动态大小的连续序列
- str：字符串切片
- `(T,U,..)`：有限序列，T 和 U 可以是不同类型
- `fn(i32) -> i32`：函数

### 变量定义与不可变
rust 中变量默认不可变，变量需要可变使用关键字 **mut**
```rust
	let foo = "a";	// 默认不可变
	let mut bar = "b"; // 变量可变
```
let 与 mut 的区别：
- 可以使用 let 关键字重复声明一个变量，每次 let 声明都是创建了一个新变量，可以改变值类型，但复用了这个名字

```rust
let a = "";
let a = 1;
let a = 1.2;
// 以上使用使用相同变量隐藏了一个变量，行为看似是可以对一个变量赋值不同的类型

// mut 关键字
let mut a = " ";
a = 2;	// 会报错，a 变量可变，但是不能改变其类型
```

### 表达式和语句
**rust 中表达式有返回值，语句没有返回值**
```rust
let y = {
	let x = 3;
	x + 1	// 因为没有分号，所以表达式会有返回值
};
println!("the value of y is: {}",y);
```

### Strings
- `&str`：指向一个存在的字符串，它可能存在 stack、heap 或者 data segment 上。
- `String`：只会分配在 heap 上。

## const 和 static
### const
使用 const 定义的变量时，它们总是会被内联。

### static
static 拥有固定的内存地址和作为一个单例存在，它们是可变的。读写 static 变量需要在 unsafe 块中
```rust
static mut BAZ: u32 = 4;
static FOO: u8 = 9;

fn main() {
	unsafe {
		println!("baz is {}",BAZ);
		BAZ = 42;
		println!("baz is now {}",BAZ);
		println!("foo is {}",FOO);
	}
}
```

### const fn
`const fn` 会在编译时被执行，传入 `const fn` 的参数必须是不可变的，不能在函数内包含带有对 heap 的操作。

```rust
const fn salt(a: u32) -> u32 {
	0xDEADBEEF ^ a
}

const CHECKSUM: u32 = salt(23);

fn main() {
	println!("{}",CHECKSUM);	
}
```

### if else
rust 中 if else 会将最后一行作为返回值，且 if else 返回类型必须一致。如果省略 else 块，则会返回一个 `()` ，rust 不允许一个变量具有两种类型。
```rust
let resutl = if 1==2 {
	"wait,what?"
} else {
	"Rust makes sense"
};

println!("You Know what? {}",result);
```

### match
```rust
	let status = req_status();
	match status {
		200 => println!("Success"),
		404 => println!("Not Found"),
		other => { // 捕获所有，如果想忽略则使用 _
			println!("Request failed with code {}",other);	
		}
	}
```

### loops
rust 支持三种循环：loop、while、for

#### loop
```rust

let mut x = 1024;
loop {
	if x < 0 {
		break;	
	}
	println!("{} more runs to go",x};
	x -= 1;
}
```

loop 支持命名，break 可以根据名字跳出指定循环
```rust

let mut x = 0;
'first: loop {
	println!("first loop");
	'second: loop {
		println!("second loop");
		if x > 5 {
			break 'first;
		}
		x += 1;
	}
}
```

#### while
```rust

let mut x = 1000;
while x > 0 {
	println!("{} more runs to go",x);
	x -= 1;
}
```

#### for
for 循环仅支持 iterators
```rust

for i in 0..10 {
	print!("{}",i);
}

println!();

for i in 0..=10 {
	print!("{}",i);
}
```

### struct
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

定义结构体方法
```rust
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
}

fn main() {
    let rect1 = Rectangle { width: 30, height: 50 };

    println!(
        "The area of the rectangle is {} square pixels.",
        rect1.area()
    );
}
```

**关联函数**
在 impl 块中定义不以 self 作为参数的函数，被称为关联函数
```rust
impl Rectangle {
    fn square(size: u32) -> Rectangle {
        Rectangle { width: size, height: size }
    }
}
```

每个结构体都允许有多个 impl 块

#### tuple struct
```rust

struct Color(u8,u8,u8);

fn main() {
	let white = Color(255,255,255);
	let red = white.0;
	let green = white.1;
	let blue = white.2;

	println!("Red: {}, green: {}, Blue: {}\n",red,green,blue);

	let orange = Color(255,165,0);
	// 类似于 js 中的解构对象
	let Color(r,g,b) = orange;
	println!("Orange R: {}, G: {}, B: {}",r,g,b);

	// 不需要的可以使用 _ 忽略
	let Color(r,_,b) = orange;
}
```

### enum
rust 枚举中的每一个成员都可以处理不同类型和数量的数据
```rust
enum IpAddr {
	V4,
	V6,
}


struct Ipv4Addr {
    // --snip--
}

struct Ipv6Addr {
    // --snip--
}

enum IpAddr {
    V4(Ipv4Addr),
    V6(Ipv6Addr),
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

> rust 没有 null，使用 option 替代
```rust
enum Option<T> {
	Some(T),
	None,
}
```


```rust

#[derive(Debug)]
enum Direction {
    N,
    E,
    S,
    W
}

enum PlayerAction {
    Move {
        direction: Direction,
        speed: u8
    },
    Wait,
    Attack(Direction)
}

fn main() {
    let simulated_player_action = PlayerAction::Move {
        direction: Direction::N,
        speed: 2,
    };

    match simulated_player_action {
        PlayerAction::Wait => println!("player wants to wait"),
        PlayerAction::Move { direction,speed } => {
            println!("player wants to move in direction {:?} with speed {}",direction,speed)
        },
        PlayerAction::Attack(direction) => {
            println!("player wants to attack direction {:?}",direction)
        }
    };
}
```

### impl blocks for struct
```rust

struct Player {
    name: String,
    iq: u8,
    friends: u8
}

impl Player {
    fn with_name(name: &str) -> Player {
        Player{
            name: name.to_string(),
            iq: 100,
            friends: 100
        }
    }

    fn get_friends(&self) -> u8 {
        self.friends
    }

    fn set_friends(&mut self,count: u8) {
        self.friends = count;
    }
}

fn main() {
    let mut player = Player::with_name("xiaoming");
    player.set_friends(23);
    println!("{} friends count: {}",player.name,player.get_friends());


    let _ = Player::get_friends(&player);
}
```

类型的三种实例方法：
- self：获取所有权，后续不可继续使用该实例
- &self：借用，仅提供读取实例
- &mut self：可变引用，允许读写该方法

### impl blocks for enum

```rust

enum PaymentMode {
    Debit,
    Credit,
    Paypal
}

fn pay_by_credit(amount: u64) {
    println!("Processing credit payment of {}",amount);
}

fn pay_by_debit(amount: u64) {
    println!("Processing debit payment of {}",amount);
}

fn paypal_redirect(amount: u64) {
    println!("Redirecting to paypal for amount {}",amount);
}

impl PaymentMode {
    fn pay(&self,amount: u64) {
        match self {
            PaymentMode::Credit => pay_by_credit(amount),
            PaymentMode::Debit => pay_by_debit(amount),
            PaymentMode::Paypal => paypal_redirect(amount)
        }
    }
}

fn get_saved_payment_mode() -> PaymentMode {
    PaymentMode::Debit
}

fn main() {
    let payment_mode = get_saved_payment_mode();
    payment_mode.pay(512);
}
```


### vectors
vectors 分配在 heap 上。
创建方式：`Vec::new` 或者 `vec![]` 宏
```rust

fn main() {
	let mut numbers_vec: Vec<u8> = Vec::new();
	numbers_vec.push(1);
	numbers_vec.push(2);

	let mut vec_with_macro = vec![1];
	vec_with_macro.push(2);
	let _ = vec_with_macro.pop();

	let message = if numbers_vec == vec_with_macro {
		"They are equal"
	} else {
		"Nah! They look different to me"
	};
	println!("{} {:?} {:?}",message,numbers_vec,vec_with_macro);
}
```


### slices
slices 表示为：`&[T]`
```rust

fn main() {
	let mut numbers: [u8;4] = [1,2,3,4];
	{
		let all: &[u8] = &numbers[..];
		println!("All of them {:?}",all);
	}

	{
		let first_two: &mut [u8] = &mut numbers[0..2];
		first_two[0] = 100;
		first_two[1] = 99;
	}
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
