---
tags: ["rust"]
date created: 2021-06-25 21:32
date modified: 2023-01-21 14:11
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
- `&str`：字符串切片
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

## 复合类型

### 字符串

- `str`：是被硬编码进二进制文件的，通常以 `&str` 形式出现
- `String`：可增长、可改变且具有所有权的 UTF-8 字符串，被分配在堆上

### 元组

元组是多种类型组合在一起形成的，元组的长度固定，元组中元素的顺序也是固定的

```rust
let tup: (i32, f64, u8) = (500, 6.4, 1);

// 解构元组
let (x,y,x) = tup;

// 访问元组
println!("{},{},{}",tup.0,tup.1,tup.2);
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

#### 元组结构体 (tuple struct）

结构体必须要有名称，但是字段可以没有名称，这种结构体被称为元组结构体

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

**枚举类型是一个类型，它包含所有可能的枚举成员，而枚举值是该类型中的具体某个成员实例**

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

// 模式匹配获取枚举值
let msg = Message::Move{x:1, y:2}
if let Message::Move{x:a, y:b} = msg  {
	println!("{},{}",a,b);
}
```

#### Option 枚举用于处理空值

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

### 数组

rust 中有两种数组，一种是长度固定的 array，第二种是可动态增长的 vector

**数组存储在栈上，vector 存储在堆上**

#### array

```rust
// 创建数组的几种方式
let a = [1,2,3,4];

let b: [i32;5] = [1,2,3,4,5];

// 某个值重复出现N次
let a = [3;5];
```

#### vectors

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

#### slices

slices 表示为：`&[T]`

```rust
fn main() {
	// 通过数组获得一个切片
	let mut numbers: [u8;4] = [1,2,3,4];
	{
		let all: &[u8] = &numbers[..];
		println!("All of them {:?}",all);
	}

	// 切片
	{
		let first_two: &mut [u8] = &mut numbers[0..2];
		first_two[0] = 100;
		first_two[1] = 99;
	}
}
```

### 流程控制

#### if

rust 中 if else 会将最后一行作为返回值，且 if else 返回类型必须一致。如果省略 else 块，则会返回一个 `()` ，rust 不允许一个变量具有两种类型。

```rust
let resutl = if 1==2 {
	"wait,what?"
} else {
	"Rust makes sense"
};

println!("You Know what? {}",result);
```

#### 循环控制

rust 支持三种循环：loop、while、for

##### loop

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

##### while

```rust
let mut x = 1000;
while x > 0 {
	println!("{} more runs to go",x);
	x -= 1;
}
```

##### for

>[!Note]

>使用 for 时，需要注意使用集合的引用形式，如果不使用引用，所有权会被转移到 for 语句块中，后面就无法再使用这个集合了

>对于实现了 copy 特征的数组，则不会产生所有权转移，而是直接进行了拷贝，因此循环后任然可以使用 arr

```rust
for i in 0..10 {
	print!("{}",i);
}

println!();

for i in 0..=10 {
	print!("{}",i);
}
```

## 模式匹配

- match 匹配需要穷举出所有可能，因此可以用 `_` 代表未列出的所有可能性
- match 的每一个分支都必须是一个表达式，且所有分支表达式最终返回值类型都必须相同
- `X|Y` 类似运算符，代表分支可以匹配 x 也可以匹配 y，只要满足一个即可
```rust
let status = req_status();
match status {
	200 => println!("Success"),
	404 => println!("Not Found"),
	// _ 相当于default
	_ => { 
		println!("Request failed with code {}",other);	
	}
}
```

模式匹配可以从模式中取出绑定的值

```rust
#[derive(Debug)]
enum UsState {
	Alabama,
	Alaska,
}

enum Coin {
	Penny,
	Nickel,
	Dime,
	Quarter(UsState),
}

let coin = Coin::Quarter(UsState::Alabama);
let num = match coin {
	Coin::Peny => 1,
	Coin::Nickel => 5,
	Coin::Dime => 10,
	Coin::Quarter(state) => {
		println!("state quarter from {:?}",state);
		25
	};
}
```

### if let 匹配

当只有一个匹配条件，且忽略其他条件时就用 `if let`

```rust
let o = Some(23);
if let Some(i) = o {
	println!("value i: {i}");
}
```

### matches! 宏

`matches!` 宏，可以将一个表达式跟模式进行匹配，然后返回 true 或 false

```rust
let foo = 'f';
assert!(matches!(foo,'A'..='Z' | 'a'..='z');
```

## 方法

rust 使用 `impl` 来定义方法

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

## 所有权和移动

### 可变性

当所有权转移时，数据的可变性可能发生改变
```rust
let immutable_box = Box::new(5u32);
// 不可变，重新赋值会报错
// *immutable_box = 4;

// 移动box，改变所有权和可变性
let mut mutable_box = immutable_box;
// 修改变量
*mutable_box = 4;
```

### 部分移动

单个变量解构，可以使用模式绑定，意味着变量的某些部分将被移动，而其他部分将保留。
这种情况下，后面不能整体使用父级变量，但是任然可以使用只引用的部分。
```rust
#[derive(Debug)]
struct Person {
	name: String,
	age: u8,
}

let person = Person {
	name: String::from("alice"),
	age: 20,
};

// name所有权被转移，age只是引用
let Person { name, ref age} = person;

// 报错，因为 name 的所有权被转移，person发生了部分借用
// println!("person: {:?}",person);

// person age 可以继续使用
// println!("person age: {}", person.age);
```

## 借用

使用 `&T` 来访问数据，而不是直接获取数据的所有权
编译器通过借用检查，静态的保证引用总是指向有效的对象，当存在引用指向一个对象时，该对象不能被销毁。

### 可变性

可变数据使用 `&mut T` 叫做可变借用，`&T` 叫做不可变借用

可变借用与不可变借用不能同时存在，且同一时间内只允许一次可变借用

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

```rust
pub trait Animal {
	eat();
	bar();
}

strcut dog {}

impl Animal for dog {
	fn eat(&self) {}
	fn bar(&self) {}
}

// 使用trait作为函数参数
pub fn action(animal: &impl Animal) {
	animal.eat();
	animal.bar();
}
```

#### 多重约束

```rust
// 多重约束
pub fn notify(item: &(impl Animal + Display)) {}

pub fn notify<T: Animal+Display>(item: &T) {}
```

#### Where 约束

当约束变多，函数签名变得复杂，使用 where 更清晰

```rust
fn some_function<T,U>(t: &T, u: &U) -> i32 
	where T: Display + Clone,
		  U: Clone + Debug
	{}
```

#### 虚类型参数

在运行时不出现，仅在编译时进行静态检查的类型参数

#### 特征对象

rust 无法通过 `impl trait` 实现不同的返回类型，因此引入了特征对象

`dyn` 关键字用在特征对象的声明上

```rust
trait Draw {
	fn draw(&self);
}

struct Button {
	width: u32,
	height: u32,
	label: String,
}

impl Button {
    fn new() -> Self {
        Self {
            width: 0,
            height: 0,
            label: "button".to_string(),
        }
    }
}

impl Draw for Button {
	fn draw(&self) {
		println!("{}", self.label);
	}
}

impl Input {
    fn new() -> Self {
        Self {
            width: 0,
            height: 0,
            label: "input".to_string(),
        }
    }
}

struct Input {
	width: u32,
	height: u32,
	label: String,
}

impl Draw for Input {
	fn draw(&self) {
		println!("{}", self.label);
	}
}

fn main() {
	let components: Vec<Box<dyn Draw>> = vec![
        Box::new(Input::new()),Box::new(Button::new()),
    ];
    for draw in components {
        draw.draw();
    }
}
```

#### 调用同名方法

```rust
trait Pilot {
	fn fly(&self);
}

trait Wizard {
	fn fly(&self);
}

struct Human;
impl Pilot for Human {
	fn fly(&self) {
		println!("pilot human fly");
	}
}

impl Wizard for Human {
	fn fly(&self) {
		println!("wizard human fly");
	}
}

impl Human {
	fn fly(&self) {
		println!("human fly");
	}
}

fn main() {
	let h = Human{};
	// 默认调用 human 上的方法
	h.fly();
	// 想调用各自 trait 上方法时，以下调用方法，只针对函数参数带有 &self
	Pilot::fly(&h);
	Wizard::fly(&h);
}
```

#### 完全限定语法

当函数参数不带 `&self` 时的调用方法

```rust
trait Animal {
	fn name() -> String;
}

struct Dog;
impl Dog {
	fn name() -> String {
		String::from("小黑")
	}
}

impl Animal for Dog {
	fn name() -> String {
		String::from("小白")
	}
}

fn main() {
	// 打印小黑
	println!("{}",Dog::name());

	// 因为函数不带&self,想打印小白的调用方法如下
	println!("{}",<Dog as Animal>::name());
}
```

### 生命周期

生命周期的主要作用是避免悬垂引用，避免程序引用了不该引用的数据

生命周期语法：以 `'` 开头，如 `'a mut i32`

```rust
// 表示first和second，至少需要有一个共同的'a生命周期，而'a生命周期，取first和second中较小的一个
fn useless<'a>(first: &'a i32, second: &'a i32) {}
```

#### 生命周期约束语法

`T: 'a`，在 T 中的所有引用都必须比 `'a` 活得更长
`T: Trait + 'a` ：T 类型必须实现 `Trait`，并且在 T 中的所有引用都必须比 `'a` 活得更长

```rust
// 'a: 'b 表示 'a 至少要活得和 'b 一样久
fn longest<'a:'b,'b>(first: &'a str, second: &'b str) -> &'b str {
    if first.len() > second.len() {
        return first;
    }

    second
}

// 约束语法，也可以使用 where 表示
fn longest<'a,'b>(first: &'a str, second: &'b str) -> &'b str 
where
	'a: 'b,
{
		.
		.
		.
}

// 'a:'b 表示'a 必须活的比 'b 久
```

#### 静态生命周期

静态生命周期的引用，存活的和程序一样久

```rust
let s: &'static str = "静态生命周期";
```

## 类型转换

### as 转换

>[!note]

>使用 as 类型转换时需要注意要转换类型的最大值，避免溢出

>强制类型转换不具有传递性，如 `e as U1 as U2` ，但不能表示 `e as U2` 是合法的

```rust
let a = 3.1 as i8;
let b = 100_i8 as i32;
let c = 'a' as u8; // 将字符转为整数

println!("{a},{b},{c}");
```

**内存地址转换为指针**

```rust
let mut values: [i32;2] = [1,2];
let p1: *mut i32 = values.as_mut_ptr();
let first_address = p1 as usize; // 将p1内存地址转换为一个整数
let second_address = first_address + 4; //将内存地址偏移一个i32的大小
let p2 = second_address as *mut i32; // 访问该地址执行的下一个整数
unsafe {
	*p2 += 1;
}

println!("{:?}",values);
```

## 错误处理

使用 `panic!()` 宏，触发 panic

### 传播错误

```rust
fn read_username_from_file() -> Result<String,io::Error> {
	let mut f = File::open("hello.txt")?; // 错误传播
	let mut s = String::new();
	f.read_to_string(&mut s)?;
	Ok(s)
}
```

`?` 相当于一个宏，简化了一下写法

```rust
let mut f = match f {
	Ok(file) => file,
	Err(e) => return Err(e),
};
```

同时 `?` 也可以自动进行类型转换

```rust
// 此处因为 std::io::Error 实现了 std::error::Error,因此 ? 可以实现转换
fn open_file() -> Result<File, Box<dyn std::error::Error>> {
	let mut f = File::open("hello.txt")?;
	Ok(f)
}
```

## 包和 package

- 项目（Package）：用来构建、测试和分享包
- 工作空间（workspace）：大型项目，可以进一步将多个包联合在一起，组织成工作空间
- 包（Crate）：一个由多个模块组成的树形结构，可以将作为第三方库进行分发，也可以生成可执行文件
- 模块（Module）：一个文件多个模块，或一个文件一个模块，模块可以被认为是真实项目中的代码组织单元

### 路径引用模块

- 绝对路径：从包根开始，路径名以包名或 `crate` 作为开头
- 相对路径：
	- self
	- super：以父模块为开始进行引用，相当于 `../a/b`

### 结构体和枚举的可见性

结构体设置为 `pub`，但它的所有字段依然是私有的

将枚举设置为 `pub` ，它的所有字段将对外可见

### use 的使用

**限制可见性语法：**

- `pub` 可见性无任何限制
- `pub(crate)` 在当前包可见
- `pub(self)` 在当前模块可见
- `pub(super)` 在父模块可见
- `pub(in <path>)` 在某个路径的模块中可见，path 必须是父模块或者祖先模块

## 函数式编程

### 闭包

闭包是一种匿名函数，不同于函数的是，它可以捕获调用者作用域的值

```rust
|param1,...| {
	语句1;
	语句2;
	返回表达式
}

// 当只有一个返回表达式时
|param1,[param2],....| 返回表达式
```

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

- FnOnce 捕获周围的变量，同时闭包会获取变量的所有权。因为获取了所有权则所以该闭包只能被调用一次。
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

### 迭代器

实现了 `IntoIterator` 特征，就可以进行迭代

- into_iter：夺走所有权
- iter：借用
- iter_mut：可变借用