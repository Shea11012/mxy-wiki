---
date created: 2021-11-05 11:40
---

## 实体关系图 Entity Relationship Diagrams

```mermaid
erDiagram
	customer ||--o{order:places
	customer {
		string name
		string custNumber
		string sector
	}
	
	order ||--|{ line-item:contains
	order {
		int orderNumber
		string deliveryAddress
	}
	
	line-item {
		string products
		int quantity
		float pricePerUnit
	}
```

| Value left | value right | description |
| ---------- | ----------- | ----------- |
| \|o | o\| |0 或1 |
| \|\| \|\| | 仅一个 |
| )o | o( | 0个或多个 |
| }\| \|{ | 一个或多个 |

```mermaid
erDiagram
	CAR ||--o{ Name-Driver: allows
	CAR {
		string registrationNumber
		string make
		string model
	}
	
	Person }|--o{ Name-Driver: is
	Person {
		string firstName
		string lastName
		int age
	}
```