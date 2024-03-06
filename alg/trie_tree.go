package alg

import (
	"bytes"
	"sort"
	"sync"
)

type TrieTreeNode struct {
	Data     byte
	Children []*TrieTreeNode
	RefCount int
}

func (t *TrieTreeNode) String() string {
	var buf = bytes.NewBuffer([]byte{})
	t.writeToString(t, buf)
	return buf.String()
}

func (t *TrieTreeNode) writeToString(node *TrieTreeNode, buf *bytes.Buffer) {
	buf.WriteByte(node.Data)
	buf.WriteByte(13)
	for _, child := range t.Children {
		if child != nil {
			t.writeToString(child, buf)
		}
	}
}

func NewTrieTree() (t *TrieTree) {
	t = new(TrieTree)
	t.Root = new(TrieTreeNode)
	return t
}

// TrieTree prefix based tree
// TrieTree 前缀树
type TrieTree struct {
	Root *TrieTreeNode
	lock sync.RWMutex
}

// Insert 将字符组插入到树上
func (t *TrieTree) Insert(words []byte) {
	t.lock.Lock()
	defer t.lock.Unlock()

	nextNode := t.Root
	// for height, w := range words {
	for _, w := range words {
		// nextNode = t.inertWord(nextNode, w, height)
		nextNode = t.inertWord(nextNode, w)
	}
}

// Match longest prefix match
// Match 最长前缀匹配
func (t *TrieTree) Match(words []byte) (matched []byte) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	matchedNode := t.Root
	// for height, w := range words {
	for _, w := range words {
		// matchedNode = t.searchNode(matchedNode, w, height)
		matchedNode = t.searchNode(matchedNode, w)
		if matchedNode != nil {
			matched = append(matched, matchedNode.Data)
		} else {
			break
		}
	}

	return matched
}

// func (t *TrieTree) inertWord(startNode *TrieTreeNode, data byte, height int) (nextNode *TrieTreeNode) {
func (t *TrieTree) inertWord(startNode *TrieTreeNode, data byte) (nextNode *TrieTreeNode) {
	var index int
	var found bool
	var childCount = len(startNode.Children)

	if childCount > 0 {
		index = sort.Search(childCount, func(i int) bool {
			return startNode.Children[i].Data >= data
		})
		if index < childCount && startNode.Children[index].Data == data {
			found = true
		}
	}

	if !found {
		newNode := &TrieTreeNode{
			Data:     data,
			RefCount: 1,
		}
		startNode.Children = append(startNode.Children, newNode)
		sort.Slice(startNode.Children, func(i, j int) bool {
			return startNode.Children[i].Data < startNode.Children[j].Data
		})
		nextNode = newNode
	} else {
		startNode.Children[index].RefCount++
		nextNode = startNode.Children[index]
	}

	return
}

// func (t *TrieTree) searchNode(startNode *TrieTreeNode, data byte, height int) (matchedNode *TrieTreeNode) {
func (t *TrieTree) searchNode(startNode *TrieTreeNode, data byte) (matchedNode *TrieTreeNode) {
	var index int
	var found bool
	var childCount = len(startNode.Children)
	if childCount > 0 {
		index = sort.Search(childCount, func(i int) bool {
			return startNode.Children[i].Data >= data
		})
		if index < childCount && startNode.Children[index].Data == data {
			found = true
		}
	}

	if found {
		matchedNode = startNode.Children[index]
	}

	return
}

// Delete delete prefix from tree
// Delete 将字符组从树上删除,输入的字符必须是从头到尾的前缀字符
func (t *TrieTree) Delete(words []byte) {
	t.lock.Lock()
	defer t.lock.Unlock()

	var matchedNodes []*TrieTreeNode
	// for height, w := range words {
	for _, w := range words {
		// matchedNode := t.searchNode(t.Root, w, height)
		matchedNode := t.searchNode(t.Root, w)
		if matchedNode != nil {
			matchedNodes = append(matchedNodes, matchedNode)
		} else {
			break
		}
	}

	for i := len(matchedNodes) - 1; i >= 0; i-- {
		if i-1 < 1 {
			t.deleteNode(t.Root, matchedNodes[i])
		} else {
			t.deleteNode(matchedNodes[i-1], matchedNodes[i])
		}
	}
}

func (t *TrieTree) deleteNode(parentNode, node *TrieTreeNode) {
	if node == t.Root {
		return
	}

	node.RefCount--

	var parentChildCount = len(parentNode.Children)
	if node.RefCount < 1 {
		index := sort.Search(parentChildCount, func(i int) bool {
			return parentNode.Children[i].Data >= node.Data
		})
		if index < parentChildCount && parentNode.Children[index].Data == node.Data {
			parentNode.Children = append(parentNode.Children[:index], parentNode.Children[index+1:]...)
		}
	}
}

// // Walk list all record combo from trie tree
// func (t *TrieTree) Walk(f func(line []byte)) {
// 	for _, child := range t.Root.Children {
// 		for i := 0; i < child.RefCount; i++ {
// 			f(t.walk(child, i+1, f))
// 		}
// 	}
// }

// // walk 从某个节点开始遍历
// func (t *TrieTree) walk(node *TrieTreeNode, childIndex int, f func(line []byte)) (lineCollect []byte) {
// 	lineCollect = append(lineCollect, node.Data)

// 	if len(node.Children) > 0 {
// 		var realChildIndex int
// 		var totalRefCount int
// 		for k, child := range node.Children {
// 			if childIndex-totalRefCount == child.RefCount {
// 				realChildIndex = k
// 				break
// 			} else {
// 				totalRefCount += child.RefCount
// 			}
// 		}
// 		lineCollect = append(lineCollect, t.walk(node.Children[realChildIndex], childIndex, f)...)
// 	}

// 	return lineCollect
// }

// WalkSuffix list all suffix by prefix, used for hot search keywords recommend
// WalkSuffix 遍历剩余后缀,用于搜热词推荐
// func (t *TrieTree) WalkSuffix(prefix []byte, f func(line []byte), count int) {

// }
