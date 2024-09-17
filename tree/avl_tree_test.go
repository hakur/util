package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAVLTreeInsert(t *testing.T) {
	tree := NewAVLTree[int]()
	tree.Insert(14)
	tree.Insert(9)
	tree.Insert(5)
	assert.Equal(t, 9, tree.Root.Value)
	assert.Equal(t, 5, tree.Root.Left.Value)
	assert.Equal(t, 14, tree.Root.Right.Value)

	tree.Insert(17)
	tree.Insert(11)
	tree.Insert(12)
	assert.Equal(t, 11, tree.Root.Value)
	assert.Equal(t, 9, tree.Root.Left.Value)
	assert.Equal(t, 5, tree.Root.Left.Left.Value)
	assert.Equal(t, 14, tree.Root.Right.Value)
	assert.Equal(t, 12, tree.Root.Right.Left.Value)
	assert.Equal(t, 17, tree.Root.Right.Right.Value)

	tree.Insert(7)
	assert.Equal(t, 11, tree.Root.Value)
	assert.Equal(t, 7, tree.Root.Left.Value)
	assert.Equal(t, 5, tree.Root.Left.Left.Value)
	assert.Equal(t, 9, tree.Root.Left.Right.Value)
	tree.Insert(19)
	tree.Insert(16)
	tree.Insert(27)
	assert.Equal(t, 11, tree.Root.Value)
	assert.Equal(t, 7, tree.Root.Left.Value)
	assert.Equal(t, 17, tree.Root.Right.Value)
}

func TestAVLTreeDelete(t *testing.T) {
	tree := NewAVLTree[int]()
	data := []int{14, 9, 5, 17, 11, 12, 7, 19, 16, 27}
	for _, v := range data {
		tree.Insert(v)
	}
	assert.Equal(t, 11, tree.Root.Value)

	tree.Delete(11)
	assert.Equal(t, 12, tree.Root.Value)
	tree.Delete(12)
	assert.Equal(t, 14, tree.Root.Value)
	assert.Equal(t, 16, tree.Root.Right.Left.Value)

	tree.Delete(5)
	tree.Delete(9)
	assert.Equal(t, 17, tree.Root.Value)

	tree.Delete(14)
	assert.Equal(t, 17, tree.Root.Value)
	assert.Equal(t, 7, tree.Root.Left.Value)
	tree.Delete(7)
	assert.Equal(t, 17, tree.Root.Value)
	assert.Equal(t, 16, tree.Root.Left.Value)

	tree.Delete(16)
	assert.Equal(t, 19, tree.Root.Value)
	assert.Equal(t, 17, tree.Root.Left.Value)

	tree.Delete(19)
	assert.Equal(t, 17, tree.Root.Value)
	assert.Equal(t, 27, tree.Root.Right.Value)

	tree.Delete(17)
	assert.Equal(t, 27, tree.Root.Value)

	tree.Delete(27)
	assert.Equal(t, true, tree.Root == nil)
}

func TestAVLTreeGetHeihgt(t *testing.T) {
	tree := NewAVLTree[int]()
	tree.Insert(14)
	tree.Insert(9)
	tree.Insert(5)
	tree.Insert(15)
	assert.Equal(t, 9, tree.Root.Value)
	assert.Equal(t, 3, tree.GetNodeHeight(tree.Root))
	assert.Equal(t, 1, tree.GetNodeHeight(tree.Root.Left))
	assert.Equal(t, 2, tree.GetNodeHeight(tree.Root.Right))
}

func TestAVLTreeDelete10000(t *testing.T) {
	tree := NewAVLTree[int]()
	for i := 0; i <= 10000; i++ {
		tree.Insert(i)
	}

	for i := 0; i <= 10000; i++ {
		tree.Delete(i)
	}

	assert.Equal(t, true, tree.Root == nil)
}

func BenchmarkAVLTreeInsert(b *testing.B) {
	tree := NewAVLTree[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(i)
	}
}

func BenchmarkAVLTreeSearch(b *testing.B) {
	tree := NewAVLTree[int]()
	for i := 0; i <= 1000000; i++ {
		tree.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(i)
	}
}

func TestAVLTreeWalk(t *testing.T) {
	tree := NewAVLTree[int]()
	data := []int{14, 9, 5, 17, 11, 12, 7, 19, 16, 27}
	for _, v := range data {
		tree.Insert(v)
	}

	var sorted []int
	tree.Walk(nil, func(node *AVLTreeNode[int]) {
		sorted = append(sorted, node.Value)
	})

	assert.Equal(t, []int{5, 7, 9, 11, 12, 14, 16, 17, 19, 27}, sorted)
}
