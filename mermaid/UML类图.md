## UML类图

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
| \*--      | composition 组合  |
| o--      | Aggregation 聚合  |
| -->      | Association 关联  |
| --       | 实线              |
| ..>      | Dependency 依赖   |
| ..\|>     | Realization 实现  |
| ..       | 虚线              |

```mermaid
classDiagram
	classA <|-- classB
	classC *-- classD
	classE o-- classF
	classG <-- classH
	classI -- classJ
	classK <.. classL
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

- 1：仅1个
- 0 .. 1：0 或 1
- 1 .. *：1 或 多
- *：多
- n：n个，n>1
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