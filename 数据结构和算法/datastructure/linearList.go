package main

import (
	"errors"
	"fmt"
)

type Elem int

type Linear struct {
	// 表示当前长度
	Length int
	// 表示总长度
	MaxLength int
	data      []Elem
}

func main() {
	l := New(8)
	l.InsertItem(3,4)
	l.InsertItem(2,8)
	l.InsertItem(1,29)
	l.InsertItem(4,20)
	l.DeleteItem(3)
	l.InsertItem(8,31)
	i,_ := l.LocateItem(20)
	fmt.Printf("locate %d\n",i)
	fmt.Printf("%v,length = %d,cap = %d,len = %d\n",l.data,l.Length,cap(l.data),len(l.data))

}

// 初始化 linear 结构
func New(maxlength int) *Linear {
	return &Linear{
		MaxLength: maxlength,
		data:      make([]Elem, maxlength),
	}
}

// 判断线性表是否为空
func (l Linear) IsEmpty() bool {
	if l.Length == 0 {
		return true
	}

	return false
}

// 清空线性表
func (l Linear) ClearList() Linear {
	return l
}

// 在第 i 个位置插入元素
func (l *Linear) InsertItem(i int,e Elem) bool {
	if i > l.MaxLength || i < 0 {
		return false
	}

	l.data[i-1] = e
	l.Length++
	return true
}

// 删除第 i 个元素
func (l *Linear) DeleteItem(i int) bool {

	if l.Length == 0 {
		return false
	}

	if i > l.MaxLength || i < 1 {
		return false
	}

	for j := i; j < l.MaxLength; j++ {
		l.data[j-1] = l.data[j]
	}

	// 如果是删除最后一个元素，不会进入 for 循环，需要把最后一位初始化
	if i == l.MaxLength {
		l.data[i-1] = 0
	}

	l.Length--
	return true
}

// 获取当前长度
func (l Linear) GetLength() int {
	return l.Length
}

// 定位元素位置
func (l *Linear) LocateItem(e Elem)(int,error) {
	for i,v := range l.data {
		if v == e {
			return i+1,nil
		}
	}

	return 0,errors.New("not found element")
}


// 将第 i 个位置元素返回
func (l Linear) getElem(i int) (Elem,bool) {
	if i<1 || i > l.MaxLength {
		return 0,false
	}
	return l.data[i],true
}
