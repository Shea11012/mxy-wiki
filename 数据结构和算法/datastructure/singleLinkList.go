package main

import "fmt"

type Linked struct {
	data string
	next *Linked
}

func main() {
	//l := Linked{
	//	data: "1",
	//	next: &Linked{
	//		data: "2",
	//		next: &Linked{
	//			data: "3",
	//			next: &Linked{
	//				data: "4",
	//				next: &Linked{
	//					data: "5",
	//					next: nil,
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//l.ListInsert(2, &Linked{
	//	data: "6",
	//	next: &Linked{
	//		data: "7",
	//		next: nil,
	//	},
	//})
	//
	//l.DeleteList(3)

	l := new(Linked)
	l.CreateTailList("mxy")
	for l != nil {
		fmt.Printf("%v",l.data)
		l = l.next
	}

}

func (l *Linked) GetElem(i int) *Linked {
	j := 1
	n := l
	for n != nil {

		if n.data != "" && j == i {
			return n
		}

		n = n.next

		j++
	}

	return nil
}

func (l *Linked) ListInsert(i int, node *Linked) bool {
	if node == nil {
		return false
	}
	r := l.GetElem(i)
	if r == nil {
		return false
	}

	next := r.next
	r.next = node
	for node != nil {
		if node.next == nil {
			node.next = next
			break
		}
		node = node.next
	}

	return true
}

func (l *Linked) DeleteList(i int) bool {
	r := l.GetElem(i-1)
	if r == nil {
		return false
	}
	r.next = r.next.next

	return true
}

// 头插法
func (l *Linked) CreateHeadList(str string) bool {
	for _ ,v := range []rune(str) {
		c := new(Linked) // 每次都会新建一个节点
		c.data = string(v)	// 给值域赋值
		c.next = l.next // 将新建节点的指针域指向头结点
		l.next = c	// 将头指针的指针域指向新建节点
	}

	return true
}

// 尾插法
func (l *Linked) CreateTailList(str string) bool {
	n := l
	for _,v := range str {
		c := new(Linked)
		c.data = string(v)
		n.next = c
		n = c
	}

	return true
}
