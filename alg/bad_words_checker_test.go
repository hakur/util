package alg

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 3.1 中文违禁词测试
func TestBadWordsCheckerChineseContains(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("敏感词")
	c.Add("违禁")

	assert.True(t, c.Contains("这是一段包含敏感词的文本"))
	assert.True(t, c.Contains("违禁内容"))
	assert.False(t, c.Contains("这是一段正常的文本"))
	assert.False(t, c.Contains(""))
}

func TestBadWordsCheckerChineseFindAll(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"敏感词", "违禁", "屏蔽"})

	matches := c.FindAll("敏感词和违禁都应该被屏蔽掉")
	assert.Len(t, matches, 3)
	assert.Equal(t, "敏感词", matches[0].Word)
	assert.Equal(t, "违禁", matches[1].Word)
	assert.Equal(t, "屏蔽", matches[2].Word)
}

func TestBadWordsCheckerChineseReplace(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("敏感词")

	result := c.Replace("这段文本包含敏感词需要屏蔽", '*')
	assert.Equal(t, "这段文本包含***需要屏蔽", result)

	// 无违禁词时返回原文本
	result = c.Replace("正常文本", '*')
	assert.Equal(t, "正常文本", result)
}

// 3.2 英文违禁词测试 + 子串匹配
func TestBadWordsCheckerEnglishContains(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil", "danger"})

	assert.True(t, c.Contains("this is bad"))
	assert.True(t, c.Contains("evil plan"))
	assert.False(t, c.Contains("good things"))
}

func TestBadWordsCheckerEnglishFindAll(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "word"})

	matches := c.FindAll("bad word")
	assert.Len(t, matches, 2)
	assert.Equal(t, "bad", matches[0].Word)
	assert.Equal(t, "word", matches[1].Word)
}

func TestBadWordsCheckerEnglishReplace(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "word"})

	result := c.Replace("bad word here", '*')
	assert.Equal(t, "*** **** here", result)
}

func TestBadWordsCheckerSubstringMatch(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("ab")

	// "ab" 在 "abc" 中应命中
	assert.True(t, c.Contains("abc"))
	assert.True(t, c.Contains("xaby"))
	assert.False(t, c.Contains("ac"))

	// FindAll 确认匹配位置
	matches := c.FindAll("abc")
	assert.Len(t, matches, 1)
	assert.Equal(t, "ab", matches[0].Word)
	assert.Equal(t, 0, matches[0].ByteStart)
	assert.Equal(t, 2, matches[0].ByteEnd)
}

func TestBadWordsCheckerMultipleSubstrings(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"ab", "bc"})

	matches := c.FindAll("abc")
	assert.Len(t, matches, 2)
}

// 3.3 边界测试
func TestBadWordsCheckerEmptyWord(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("") // 空字符串应被忽略
	c.Add("abc")

	assert.True(t, c.Contains("abc"))
	assert.False(t, c.Contains("")) // 空文本
}

func TestBadWordsCheckerDuplicateAdd(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("abc")
	c.Add("abc") // 重复添加，幂等
	c.Add("abc")

	matches := c.FindAll("abc abc abc")
	assert.Len(t, matches, 3)
}

func TestBadWordsCheckerNoWords(t *testing.T) {
	c := NewBadWordsChecker()

	assert.False(t, c.Contains("anything"))
	assert.Nil(t, c.FindAll("anything"))
	assert.Equal(t, "anything", c.Replace("anything", '*'))
}

func TestBadWordsCheckerOverlappingReplace(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"ab", "bc"})

	// 重叠匹配应合并后替换
	result := c.Replace("abc", '*')
	assert.Equal(t, "***", result)
}

func TestBadWordsCheckerACAutomaton(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"he", "she", "his", "hers"})

	// "ushe" 应同时匹配 "she" 和 "he"（通过失败链接链）
	matches := c.FindAll("ushe")
	assert.Len(t, matches, 2)

	// Contains 也应该检测到
	assert.True(t, c.Contains("ushe"))
	assert.True(t, c.Contains("this is his book"))
	assert.True(t, c.Contains("hers"))
	assert.False(t, c.Contains("nothing wrong"))
}

func TestBadWordsCheckerBuildOnlyWhenDirty(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("test")

	// 首次 Contains 触发 Build
	assert.True(t, c.Contains("test sentence"))

	// 连续调用不重复 Build
	assert.True(t, c.Contains("another test"))
	assert.False(t, c.Contains("clean text"))
}

// 3.4 并发安全测试
func TestBadWordsCheckerConcurrentReadWrite(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil"})

	var wg sync.WaitGroup

	// 并发读
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Contains("this is bad")
				c.FindAll("evil plan")
				c.Replace("bad evil", '*')
			}
		}()
	}

	// 并发写
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				c.Add("new")
				c.AddWords([]string{"extra", "more"})
			}
		}(i)
	}

	wg.Wait()
}

func TestBadWordsCheckerConcurrentBuild(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("test")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 并发触发首次检查，应只有一个执行 Build
			assert.True(t, c.Contains("test"))
		}()
	}
	wg.Wait()
}

// 3.5 基准测试
func BenchmarkBadWordsCheckerContains(b *testing.B) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil", "danger", "敏感词", "违禁", "屏蔽"})
	c.Build()

	text := "这是一段包含敏感词和bad内容的测试文本用于基准测试"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Contains(text)
	}
}

func BenchmarkBadWordsCheckerFindAll(b *testing.B) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil", "danger", "敏感词", "违禁", "屏蔽"})
	c.Build()

	text := "这是一段包含敏感词和bad内容的测试文本"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.FindAll(text)
	}
}

func BenchmarkBadWordsCheckerReplace(b *testing.B) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil", "danger", "敏感词", "违禁", "屏蔽"})
	c.Build()

	text := "这是一段包含敏感词和bad内容的测试文本"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Replace(text, '*')
	}
}
