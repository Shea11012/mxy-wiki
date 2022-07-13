---
date created: 2021-11-07 13:16
date modified: 2021-12-26 03:05
title: UML类图
---
## UML 类图

```mermaid
classDiagram
	class BankAccount{
		+string owner
		+int64 balance
		+Deposit(amount) bool
		+withdraw(amount) int
	}
```

### 泛型

```mermaid
classDiagram
	class Square~Shape~{
		int id
		List~int~ position
		SetPoints(List~int~ points)
		GetPoints() List~int~
	}
	
	Square: List~string~ messages
	Square: +SetMessage(List~string~ messages)
	Square: +GetMessages() List~string~
```

### 关联关系

| **Type** | **Description** |
| -------- | --------------- |
| <--      | Inheritance 继承  |
| \*--      | composition 组合；一个对象由一个对象或多个其他对象实例构成  |
| o--      | Aggregation 聚合；聚合关系中一个对象拥有一组其他对象  |
| -->      | Association 关联  |
| --       | 实线              |
| ..>      | Dependency 依赖   |
| ..\|>     | Realization 实现  |
| ..       | 虚线              |

```mermaid
classDiagram
	%% Cat 和 Dog 是 Animal 的具体实现
	class Animal
	Animal: +makeSound()
	class Cat
	Cat: +makeSound()
	class Dog
	Dog: +makeSound()
	Animal <|-- Cat
	Animal <|-- Dog
	
	大学 *--> 院系
	%% 院系包含教授
	院系 o--> 教授
	%% 教授关联学生
	教授 --> 学生
	classI -- classJ
	%% 教授依赖课程
	教授 <.. 课程
	classM <|.. classN
	classO .. classP
```

### 双向关联

```mermaid
classDiagram
	Animal <|--|> Zebra
```

### 多重关联

可选的基数

- 1：仅 1 个
- 0 .. 1：0 或 1
- 1 .. *：1 或 多
- *：多
- n：n 个，n>1
- 0 .. n：0 或 n

语法：

```
classA "基数" 箭头 "基数" classB: text
```

```mermaid
classDiagram
	Customer "1" --> "*" Ticket
	Student "1" --> "1..*" Course
	Galaxy --> "*" Star: Contains
```

### 类注解

- `<<Interface>>`
- `<<abstract>>`
- `<<Service>>`
- `<<enumeration>>`

```mermaid
classDiagram
	class Shape
	<<interface>> Shape
	
	class Person{
		<<interface>>
		int sex
		string name
		Hello() string
	}
```