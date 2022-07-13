---
date created: 2021-11-07 13:16
date modified: 2021-11-07 13:16
title: 用户行程图 User Journey Diagram
---
## 用户行程图 User Journey Diagram
语法：
`task name: score: [actor,actor,...]`

```mermaid
journey
	title my working day
	section go to work
		make tea: 5: me
		go upstairs: 3: me
		do work: 1: me,cat
	section Go home
		go downstairs: 5: me
		sit down: 5: me
```