package main

import "fmt"

type StaticList struct {
	data string
	cur  int
}

const MAXSIZE = 10

func main() {
	var a [MAXSIZE]StaticList

	for i := 0; i < MAXSIZE; i++ {
		a[i].cur = i + 1
	}

	a[MAXSIZE-1].cur = 0
	b := a[:]
	insertElem(b, 1, "A")
	insertElem(b, 2, "B")
	insertElem(b, 2, "C")
	insertElem(b, 1, "D")
	insertElem(b, 3, "E")
	DeleteElem(b,1)
	for j, v := range a {
		fmt.Printf("index:%d cur:%d data:%s\n", j, v.cur, v.data)
	}

	fmt.Printf("length:%d\n", ListLength(b))
}

// 在第 i 个位置前插入元素,必须是按照顺序插入
func insertElem(arr []StaticList, i int, elem string) bool {

	if i < 1 || i > ListLength(arr)+1 {
		return false
	}

	r := mallocSll(arr) // 获取空闲的分量浮标
	k := MAXSIZE - 1     // 获取最后一个数组的位置

	if r != 0 {
		arr[r].data = elem // 在空闲浮标上放置数据

		for j := 1; j <= i-1; j++ { // i 指定循环次数，获取对应的浮标值
			k = arr[k].cur
		}
		arr[r].cur = arr[k].cur // 将 k 对应的浮标赋值到 r 对应的浮标下
		arr[k].cur = r          // 将 r 值赋值到对应的 k 浮标下
	}
	return true
}

// 删除第 i 个元素
func DeleteElem(l []StaticList, i int) bool {
	if i < 1 || i > ListLength(l) {
		return false
	}

	k := MAXSIZE - 1

	for j := 1; j <= i-1; j++ {
		k = l[k].cur
	}

	f := l[k].cur // 找到要释放的浮标
	l[f].data = ""

	l[k].cur = l[f].cur // 将 f 的浮标给 k 的浮标
	Freesll(l, f)	// 释放 f，进入空闲浮标

	return true
}

// 获取空闲的分量浮标
func mallocSll(arr []StaticList) int {
	i := arr[0].cur

	if arr[0].cur != 0 {
		arr[0].cur = arr[i].cur // 定位到空闲分量，取出空闲分量的浮标，放进后备浮标中
	}

	return i
}

// 将删除的分量置入空闲浮标中
func Freesll(l []StaticList, k int) bool {
	l[k].cur = l[0].cur
	l[0].cur = k
	return true
}

// 统计静态链表的长度
func ListLength(a []StaticList) int {
	j := 0
	k := a[MAXSIZE-1].cur
	for k != 0 {
		k = a[k].cur
		j++
	}

	return j
}
