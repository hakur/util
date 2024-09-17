package tree

import (
	"math"
	"sort"
)

type BTreeValueType interface {
	string | int | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64
}

type BTreeNode[VT BTreeValueType] struct {
	Parent   *BTreeNode[VT]
	Values   []VT
	Children []*BTreeNode[VT]
}

// InsertValue insert new value on this node and sort values
// InsertValue 在当前节点插入新的值并重新排序
func (t *BTreeNode[VT]) InsertValue(value VT) {
	t.Values = append(t.Values, value)
	sort.Slice(t.Values, func(i, j int) bool { return t.Values[i] < t.Values[j] })
}

// InsertChild insert new child on this node and sort children
// InsertChild 在当前节点插入新的子节点并重新排序
func (t *BTreeNode[VT]) InsertChild(child *BTreeNode[VT]) {
	t.Children = append(t.Children, child)
	child.Parent = t
	if len(child.Values) > 0 && len(t.Children) > 0 {
		sort.Slice(t.Children, func(i, j int) bool { return t.Children[i].Values[0] < t.Children[j].Values[0] })
	}
}

// DownMerge do btree overflow on this node
// DownMerge 在当前节点执行B树上溢(分裂)
func (t *BTreeNode[VT]) RiseSplit() (newParent *BTreeNode[VT]) {
	var middleKey = int(math.Ceil(float64(len(t.Values))/2)) - 1
	if t.Parent == nil {
		t.Parent = new(BTreeNode[VT])
		t.Parent.InsertChild(t)
	}

	var leftChildren []*BTreeNode[VT]
	var rightChildren []*BTreeNode[VT]

	leftValues := t.Values[:middleKey]
	rightValues := t.Values[middleKey+1:]

	for k, child := range t.Children {
		if k > middleKey {
			rightChildren = append(rightChildren, child)
		} else if k <= middleKey {
			leftChildren = append(leftChildren, child)
		}
	}

	t.Parent.InsertValue(t.Values[middleKey])
	t.Values = leftValues
	t.Children = leftChildren

	var newChild = new(BTreeNode[VT])
	newChild.Values = append(newChild.Values, rightValues...)
	for _, v := range rightChildren {
		v.Parent = newChild
		newChild.Children = rightChildren
	}

	t.Parent.InsertChild(newChild)

	return t.Parent
}

// DownMerge do btree underflow on this node
// DownMerge 在当前节点执行B树下溢
func (t *BTreeNode[VT]) DownMerge(minValuesCount int) (effectNode *BTreeNode[VT]) {
	// 检查左右兄弟节点是否能借用
	if t.Parent == nil { // 如果是根节点
		// 检查根节点是不是空了，空了就合并子节点
		return
	}
	var parentChildIndex int
	var found bool
	if len(t.Values) < 1 {
		for k, child := range t.Parent.Children {
			if child == t {
				parentChildIndex = k
				found = true
				break
			}
		}
	} else {
		parentChildIndex, found = sort.Find(len(t.Parent.Children), func(i int) int {
			if t.Parent.Children[i].Values[0] == t.Values[0] {
				return 0
			} else if t.Parent.Children[i].Values[0] < t.Values[0] {
				return 1
			} else {
				return -1
			}
		})
	}
	if !found {
		return
	}

	var borrowedSiblingIndex int
	var canBorrowValue bool
	if parentChildIndex == 0 { // 最左节点只能向右邻节点借value
		borrowedSiblingIndex = parentChildIndex + 1
		if len(t.Parent.Children[borrowedSiblingIndex].Values) > minValuesCount { // 检查右邻节点够不够借
			canBorrowValue = true
		}
	} else if parentChildIndex == len(t.Parent.Values) { // 最右节点只能向左邻节点借value
		borrowedSiblingIndex = parentChildIndex - 1
		if len(t.Parent.Children[borrowedSiblingIndex].Values) > minValuesCount { // 检查左邻节点够不够借
			canBorrowValue = true
		}
	} else { // 中间节点
		borrowedSiblingIndex = parentChildIndex + 1
		if len(t.Parent.Children[borrowedSiblingIndex].Values) > minValuesCount { // 检查左节邻点够不够借
			canBorrowValue = true
		}
		if !canBorrowValue {
			borrowedSiblingIndex = parentChildIndex - 1
			if len(t.Parent.Children[borrowedSiblingIndex].Values) > minValuesCount { // 检查右邻节点够不够借
				canBorrowValue = true
			}
		}
	}

	borrowNode := t.Parent.Children[borrowedSiblingIndex]
	if canBorrowValue {
		if borrowedSiblingIndex < parentChildIndex {
			// 借用左邻节点的value来平衡自身时，将会借走左邻节点的最后一个value，并将左邻节点的第一个child插入到自身children中的第一个位置
			borrowIndex := len(borrowNode.Values) - 1
			riseUpValue := borrowNode.Values[borrowIndex]
			borrowedValue := t.Parent.Values[borrowedSiblingIndex]

			t.Parent.Values[borrowedSiblingIndex] = riseUpValue
			borrowNode.Values = append(borrowNode.Values[:borrowIndex], borrowNode.Values[borrowIndex+1:]...)
			t.Values = append([]VT{borrowedValue}, t.Values...)
		} else if borrowedSiblingIndex > parentChildIndex {
			// 借用右邻节点的value来平衡自身时，将会借走右邻节点的第一个value，并将右邻节点的第一个child合并到自身的children中
			borrowIndex := 0
			riseUpValue := borrowNode.Values[borrowIndex]
			borrowedValue := t.Parent.Values[borrowedSiblingIndex-1]

			t.Parent.Values[borrowedSiblingIndex-1] = riseUpValue
			borrowNode.Values = append(borrowNode.Values[:borrowIndex], borrowNode.Values[borrowIndex+1:]...)
			t.Values = append(t.Values, borrowedValue)
		}
		return borrowNode
	}

	// 左右相邻节点的values都不够借用时，需要进行合并操作
	if borrowedSiblingIndex < parentChildIndex { // 节点与左邻节点进行合并时，先将左邻节点对应的上级节点的value下移到左邻节点，自身的values合并到左邻节点，然后删除自身节点，然后检查上级节点是否需要执行下溢
		downMergeValue := t.Parent.Values[parentChildIndex-1]
		t.Parent.Values = t.Parent.Values[:len(t.Parent.Values)-1]

		borrowNode.Values = append(borrowNode.Values, downMergeValue)                                               // 将上一级节点的value下沉给当前节点
		borrowNode.Children = append(borrowNode.Children, t.Children...)                                            // 将子节点迁移到左邻节点
		t.Parent.Children = append(t.Parent.Children[:parentChildIndex], t.Parent.Children[parentChildIndex+1:]...) // 从父级中删除自身
	} else { // 节点与右邻节点进行合并时，现将上级节点的value移动到节点自身的values，再将右邻节点的values合并到自身的values，从上级节点删除右邻节点，然后检查上级节点是否需要执行下溢
		downMergeValue := t.Parent.Values[borrowedSiblingIndex-1]
		// t.Parent.Values = t.Parent.Values[:len(t.Parent.Values)-1]
		t.Parent.Values = append(t.Parent.Values[:borrowedSiblingIndex-1], t.Parent.Values[borrowedSiblingIndex:]...)
		t.Values = append(t.Values, downMergeValue) // 将上一级节点的value下沉给当前节点

		t.Values = append(t.Values, borrowNode.Values...)
		t.Children = append(t.Children, borrowNode.Children...)
		t.Parent.Children = append(t.Parent.Children[:borrowedSiblingIndex], t.Parent.Children[borrowedSiblingIndex+1:]...)
	}

	return t.Parent
}

type BTreeOpts struct {
	// Height height of btree
	// Height b树的高度
	Height int
	// Storage storage instance that use access keyword point data
	// Storage 用于访问关键词指向数据的存储实例
	// Storage Storage
	// AllowRepeatKey allow insert repeat key, use overflow page save repeated keys, that can keep b tree search speed fast
	// AllowRepeatKey 重复值写入，使用溢出页来保存重复值以保证B书的查找效率，https://blog.csdn.net/SStringss/article/details/123925478
	// AllowRepeatKey bool
}

// NewBTree create new btree, current not support insert reapeated value, will use overflow page resolve this problem, btree was designed for less disk io, not for fast search, use avl tree for more fast search than btree
// NewBTree 新建b树,当前不支持插入重复值，计划使用溢出页来解决重复值插入问题，b树被设计为更少的磁盘IO，并非为了快速搜索，使用avl树来获得比b树更快的搜索
func NewBTree[VT BTreeValueType](opts *BTreeOpts) *BTree[VT] {
	t := new(BTree[VT])
	t.Opts = opts
	t.Root = new(BTreeNode[VT])
	return t
}

type BTree[VT BTreeValueType] struct {
	Root *BTreeNode[VT]
	Opts *BTreeOpts
}

// Insert insert data to btree
func (t *BTree[VT]) Insert(value VT) (err error) {
	var node *BTreeNode[VT] = t.findSlotNode(t.Root, value)
	node.InsertValue(value)
	t.checkRiseSplit(node)
	return
}

func (t *BTree[VT]) checkRiseSplit(node *BTreeNode[VT]) {
	if len(node.Values) > t.Opts.Height-1 {
		var newParent = node.RiseSplit()
		if t.Root == node {
			t.Root = newParent
		}
		t.checkRiseSplit(newParent)
	}
}

func (t *BTree[VT]) findSlotNode(node *BTreeNode[VT], value VT) (targetNode *BTreeNode[VT]) {
	if len(node.Children) < 1 { // 只有叶子节点才有资格插入数据，并完成向上分裂的过程
		return node
	}

	// 枝干节点不进行数据插入，无限往下搜索到合适的叶子节点
	var childIndex int = sort.Search(len(node.Values), func(i int) bool { return value < node.Values[i] })
	if len(node.Children)-1 < childIndex {
		var newChild = new(BTreeNode[VT])
		node.InsertChild(newChild)
	}

	return t.findSlotNode(node.Children[childIndex], value)
}

func (t *BTree[VT]) Delete(value VT) (err error) {
	// 删除枝干节点的value操作都要转换为删除叶节点value的操作，选择自身开始的枝干递归到叶节点上大于自身value的的第一个value来替代自身value，然后将value从叶子节点的values中删除，最后检查被操作的叶子节点是否需要执行下溢
	node, valueIndex := t.equalSearchNode(t.Root, value)
	if node == nil {
		return
	}

	if len(node.Children) < 1 { // 如果是叶子节点直接执行删除即可
		node.Values = append(node.Values[:valueIndex], node.Values[valueIndex+1:]...)
		if node.Parent != nil {
			t.checkDownMerge(node)
		}
		return
	}
	// b树的枝干节点至少有m/2个子节点，b树根节点必须有2个子节点
	leftNode := t.getFirstLeftLeafNode(node.Children[valueIndex+1])
	if leftNode == nil {
		return
	}
	node.Values[valueIndex] = leftNode.Values[0]
	leftNode.Values = leftNode.Values[1:]
	t.checkDownMerge(leftNode)
	return
}

func (t *BTree[VT]) getFirstLeftLeafNode(node *BTreeNode[VT]) (targetNode *BTreeNode[VT]) {
	if len(node.Children) < 1 {
		return node
	}
	return t.getFirstLeftLeafNode(node.Children[0])
}

func (t *BTree[VT]) checkDownMerge(node *BTreeNode[VT]) {
	if node == nil {
		return
	}

	var minValuesCount int = int(math.Ceil(float64(t.Opts.Height)/2)) - 1
	if len(node.Values) < minValuesCount {
		effectNode := node.DownMerge(minValuesCount)
		if effectNode == t.Root {
			// root 可能的值可能被全干掉了，可能执行了左合并或者执行了右合并，但children操作后应该只会生下一个
			if len(t.Root.Values) < 1 {
				t.Root = node.Parent.Children[0] // TODO 清理一下节点链接关系
				t.Root.Parent = nil
				for _, child := range t.Root.Children {
					child.Parent = t.Root
				}
			}
		} else {
			t.checkDownMerge(effectNode)
		}
	}
}

// Walk LNR InOrder list element from small to big
// Walk 中序遍历，由小到大遍历
func (t *BTree[VT]) Walk(startNode *BTreeNode[VT], callback func(node *BTreeNode[VT], value VT)) {
	if startNode == nil {
		startNode = t.Root
	}

	t.walk(startNode, callback)
}

func (t *BTree[VT]) walk(startNode *BTreeNode[VT], callback func(node *BTreeNode[VT], value VT)) {
	if startNode == nil {
		return
	}

	for k, v := range startNode.Values {
		if len(startNode.Children) > k {
			t.walk(startNode.Children[k], callback)
		}
		callback(startNode, v)
	}
	if len(startNode.Children) > 0 {
		t.walk(startNode.Children[len(startNode.Values)], callback)
	}
}

func (t *BTree[VT]) Search(value VT) (valueContent VT, found bool) {
	node, valueIndex := t.equalSearchNode(t.Root, value)
	if node == nil {
		return valueContent, false
	}
	valueContent = node.Values[valueIndex]
	found = true
	return
}

func (t *BTree[VT]) equalSearchNode(node *BTreeNode[VT], value VT) (targetNode *BTreeNode[VT], valueIndex int) {
	var found bool
	valueIndex, found = sort.Find(len(node.Values), func(i int) int {
		if node.Values[i] == value {
			return 0
		} else if node.Values[i] < value {
			return 1
		} else {
			return -1
		}
	})

	if found {
		return node, valueIndex
	}

	valueIndex = sort.Search(len(node.Values), func(i int) bool { return value < node.Values[i] })
	if len(node.Children) < valueIndex {
		return nil, valueIndex
	}

	return t.equalSearchNode(node.Children[valueIndex], value)
}
