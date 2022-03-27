trie 树的简单实现
```go
package main

import "fmt"

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[rune]*TrieNode),
		End:      false,
	}
}

type TrieFilter struct {
	Root        *TrieNode
	Placeholder rune
}

func (t *TrieFilter) Add(kw string) {
	if len(kw) < 1 {
		return
	}
	chars := []rune(kw)

	node := t.Root

	for i := 0; i < len(chars); i++ {
		if _, ok := node.Children[chars[i]]; !ok {
			node.Children[chars[i]] = NewTrieNode()
		}
		node = node.Children[chars[i]]
	}
	node.End = true
}

func (t *TrieFilter) Replace(text string) string {
	chars := []rune(text)
	length := len(chars)

	for i := 0; i < length; i++ {
		node := t.Root
		// 不存在则跳过
		node, ok := node.Children[chars[i]]
		if !ok {
			continue
		}

		var j int
		for j = i + 1; ; j++ {
			node, ok = node.Children[chars[j]]
			if !ok {
				j = -1
				break
			}

			// node.End为true则表示是一个完整的敏感词
			if node.End {
				break
			}
		}

		if j > 0 {
			replace(chars, i, j, t.Placeholder)
			// 替换完成后，将i移至j+1上
			i = j+1
		}
	}

	return string(chars)
}

func replace(text []rune, start, end int, placeHolder rune) {
	for i := start; i <= end; i++ {
		text[i] = placeHolder
	}
}

func NewTrieFilter() *TrieFilter {
	return &TrieFilter{
		Root:        NewTrieNode(),
		Placeholder: '*',
	}
}

func main() {
	text := "狗粉丝在赌博场所说黑话"

	filter := NewTrieFilter()
	filter.Add("赌博")
	filter.Add("黑话")
	filter.Add("狗粉丝")

	fmt.Println(filter.Replace(text))
	fmt.Println(filter.Replace("狗这是一段正常的文字"))
	fmt.Println(filter.Replace("狗赌赌博"))

}

```