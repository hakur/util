package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTreeInsert4(t *testing.T) {
	var err error
	var data = []int{1, 5, 7, 4, 16, 35, 24, 42, 21, 17, 18}

	tree := NewBTree[int](&BTreeOpts{Height: 4})
	for _, v := range data {
		err = tree.Insert(v)
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, []int{7}, tree.Root.Values)
	assert.Equal(t, []int{17, 24}, tree.Root.Children[1].Values)
}

func TestBTreeInsert5(t *testing.T) {
	var err error
	var data = []int{22, 5, 11, 36, 45, 1, 3, 6, 8, 9, 13, 15, 30, 35, 40, 42, 47, 48, 50, 56}

	tree := NewBTree[int](&BTreeOpts{Height: 5})
	for _, v := range data {
		err = tree.Insert(v)
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, []int{22}, tree.Root.Values)
	assert.Equal(t, []int{36, 45}, tree.Root.Children[1].Values)
	assert.Equal(t, []int{47, 48, 50, 56}, tree.Root.Children[1].Children[2].Values)
}

func BenchmarkBTreeInsert5(b *testing.B) {
	var data = []int{22, 5, 11, 36, 45, 1, 3, 6, 8, 9, 13, 15, 30, 35, 40, 42, 47, 48, 50, 56}

	tree := NewBTree[int](&BTreeOpts{Height: 5})
	for _, v := range data {
		tree.Insert(v)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Search(56)
	}
}

func BenchmarkBTreeSearch(b *testing.B) {
	tree := NewBTree[int](&BTreeOpts{Height: 5})
	for i := 0; i < 1000000; i++ {
		tree.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(i)
	}
}

func BenchmarkBTreeInsert(b *testing.B) {
	tree := NewBTree[int](&BTreeOpts{Height: 5})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(i)
	}
}

func TestBTreeDelete(t *testing.T) {
	var err error
	var data = []int{1, 4, 5, 7, 16, 35, 24, 42, 21, 17, 18}

	tree := NewBTree[int](&BTreeOpts{Height: 4})
	for _, v := range data {
		err = tree.Insert(v)
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, []int{7}, tree.Root.Values)
	assert.Equal(t, []int{17, 24}, tree.Root.Children[1].Values)

	tree.Delete(35)
	assert.Equal(t, []int{18, 21}, tree.Root.Children[1].Children[1].Values)
	tree.Delete(42)
	assert.Equal(t, []int{18}, tree.Root.Children[1].Children[1].Values)
	tree.Delete(24)
	assert.Equal(t, []int{18, 21}, tree.Root.Children[1].Children[1].Values)

	tree.Delete(7)
	assert.Equal(t, []int{16}, tree.Root.Values)
	assert.Equal(t, []int{17}, tree.Root.Children[1].Children[0].Values)
	assert.Equal(t, []int{18}, tree.Root.Children[1].Values)
	assert.Equal(t, []int{21}, tree.Root.Children[1].Children[1].Values)

	tree.Delete(17)
	assert.Equal(t, []int{4, 16}, tree.Root.Values)
	assert.Equal(t, []int{1}, tree.Root.Children[0].Values)
	assert.Equal(t, []int{5}, tree.Root.Children[1].Values)
	assert.Equal(t, []int{18, 21}, tree.Root.Children[2].Values)
	// println(tree.Root, tree.Root.Children[0].Parent, tree.Root.Children[1].Parent, tree.Root.Children[2].Parent)

	tree.Delete(1)
	assert.Equal(t, []int{16}, tree.Root.Values)
	assert.Equal(t, []int{4, 5}, tree.Root.Children[0].Values)
	// println(tree.Root, tree.Root.Children[0].Parent, tree.Root.Children[1].Parent)

	tree.Delete(16)
	assert.Equal(t, []int{18}, tree.Root.Values)
	assert.Equal(t, []int{21}, tree.Root.Children[1].Values)
	// println(tree.Root, tree.Root.Children[0].Parent, tree.Root.Children[1].Parent)

	tree.Delete(21)
	assert.Equal(t, []int{5}, tree.Root.Values)
	assert.Equal(t, []int{4}, tree.Root.Children[0].Values)
	assert.Equal(t, []int{18}, tree.Root.Children[1].Values)

	tree.Delete(18)
	assert.Equal(t, []int{4, 5}, tree.Root.Values)
	assert.Equal(t, 0, len(tree.Root.Children))

	tree.Delete(4)
	assert.Equal(t, []int{5}, tree.Root.Values)

	tree.Delete(5)
	assert.Equal(t, []int{}, tree.Root.Values)
}

func TestBTreeWalk(t *testing.T) {
	var err error
	var data = []int{1, 4, 5, 7, 16, 35, 24, 42, 21, 17, 18}

	tree := NewBTree[int](&BTreeOpts{Height: 4})
	for _, v := range data {
		err = tree.Insert(v)
		assert.Equal(t, nil, err)
	}

	var arr []int
	tree.Walk(nil, func(node *BTreeNode[int], value int) {
		arr = append(arr, value)
	})

	assert.Equal(t, []int{1, 4, 5, 7, 16, 17, 18, 21, 24, 35, 42}, arr)
}
