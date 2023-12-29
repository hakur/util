package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieTreeInsert(t *testing.T) {
	tree := NewTrieTree()
	var matched []byte
	tree.Insert([]byte("abcdefg"))

	matched = tree.Match([]byte("abcdefg"))
	assert.Equal(t, []byte("abcdefg"), matched)

	matched = tree.Match([]byte("abc123"))
	assert.Equal(t, []byte("abc"), matched)

	tree.Insert([]byte("1.2.4/24"))
	tree.Insert([]byte("1.2.3/24"))

	tree.Insert([]byte("192.168.4.0/24"))
	tree.Insert([]byte("192.168.2.0/24"))
	tree.Insert([]byte("192.168.3.0/24"))
	// buf, err := json.Marshal(tree.Root)
	// fmt.Println(string(buf), err)
	matched = tree.Match([]byte("192.168.3.6"))
	assert.Equal(t, []byte("192.168.3."), matched)
	matched = tree.Match([]byte("192.168.2.6"))
	assert.Equal(t, []byte("192.168.2."), matched)
	matched = tree.Match([]byte("192.168.4.6"))
	assert.Equal(t, []byte("192.168.4."), matched)
}

func TestTrieTreeDelete(t *testing.T) {
	tree := NewTrieTree()

	tree.Insert([]byte("abcdefg"))
	tree.Delete([]byte("abcdefg"))
	matched := tree.Match([]byte("abcdefg"))
	assert.Equal(t, []byte(nil), matched)

	tree.Insert([]byte("abcdefg"))
	tree.Insert([]byte("abcd"))
	tree.Delete([]byte("abcd"))
	// buf, err := json.Marshal(tree.Root)
	// fmt.Println(string(buf), err)

	matched = tree.Match([]byte("abcdefg"))
	assert.Equal(t, []byte("abcdefg"), matched)
	matched = tree.Match([]byte("abc123"))
	assert.Equal(t, []byte("abc"), matched)

	tree.Delete([]byte("abcdefg"))
	matched = tree.Match([]byte("abcdefg"))
	assert.Equal(t, []byte(nil), matched)
}

// func TestTrieTreeWalk(t *testing.T) {
// 	tree := NewTrieTree()
// 	tree.Insert([]byte("1.2.4/24"))
// 	tree.Insert([]byte("1.2.3/24"))
// 	tree.Insert([]byte("1.2.5/24"))
// 	tree.Insert([]byte("1.3.5/24"))
// 	tree.Insert([]byte("1.3.6/24"))
// 	tree.Insert([]byte("1.4.6/24"))
// 	// tree.Insert([]byte("1.4.7/24"))
// 	// tree.Insert([]byte("1.4.8/24"))

// 	// buf, err := json.Marshal(tree.Root)
// 	// fmt.Println(string(buf), err)

// 	tree.Walk(func(line []byte) {
// 		println("--", string(line))
// 	})
// }

func BenchmarkTrieTreeMatch(b *testing.B) {
	tree := NewTrieTree()

	tree.Insert([]byte("abcdefg"))
	tree.Insert([]byte("1.2.4/24"))
	tree.Insert([]byte("1.2.3/24"))
	tree.Insert([]byte("192.168.4.0/24"))
	tree.Insert([]byte("192.168.2.0/24"))
	tree.Insert([]byte("192.168.3.0/24"))

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		tree.Match([]byte("192.168.3.6"))
	}
}
