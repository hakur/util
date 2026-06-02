package alg

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 3.1 中文违禁词测试（回归：应保持子串匹配）
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

	result = c.Replace("正常文本", '*')
	assert.Equal(t, "正常文本", result)
}

// 3.2 英文违禁词测试（词边界匹配）
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

// 英文词边界测试
func TestBadWordsCheckerWordBoundary(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("sm")

	// 词边界命中
	assert.True(t, c.Contains("sm test"))
	assert.True(t, c.Contains("test sm"))
	assert.True(t, c.Contains(".sm!"))
	assert.True(t, c.Contains("sm"))

	// 单词内部不命中
	assert.False(t, c.Contains("smile"))
	assert.False(t, c.Contains("small"))
	assert.False(t, c.Contains("passme"))

	// FindAll 词边界
	matches := c.FindAll("sm test sm")
	assert.Len(t, matches, 2)
}

func TestBadWordsCheckerWordBoundaryEdges(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("bad")

	// 首部边界
	assert.True(t, c.Contains("bad apple"))
	// 尾部边界
	assert.True(t, c.Contains("very bad"))
	// 标点边界
	assert.True(t, c.Contains("(bad)"))
	// 单词内部
	assert.False(t, c.Contains("badger"))
	assert.False(t, c.Contains("embad"))
}

// 中文子串回归：不应受英文词边界影响
func TestBadWordsCheckerChineseSubstringRegression(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("敏感词")

	// 中文在文本中任意位置都应命中
	assert.True(t, c.Contains("这是敏感词文本"))
	assert.True(t, c.Contains("前缀敏感词"))
	assert.True(t, c.Contains("敏感词后缀"))
}

// 更新：全 ASCII 违禁词不再命中单词内部
func TestBadWordsCheckerSubstringMatch(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("ab")

	// "ab" 在 "abc" 中 → ASCII 词边界检查失败（后面是 'c'）
	assert.False(t, c.Contains("abc"))
	assert.False(t, c.Contains("xaby"))

	// 词边界处命中
	assert.True(t, c.Contains("ab cd"))
	assert.True(t, c.Contains("ab"))
}

// 更新：重叠 ASCII 违禁词在无边界时不命中
func TestBadWordsCheckerMultipleSubstrings(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"ab", "bc"})

	// "abc" 中无词边界，两个都不命中
	matches := c.FindAll("abc")
	assert.Len(t, matches, 0)

	// 有边界时各自命中
	matches = c.FindAll("ab bc")
	assert.Len(t, matches, 2)
}

// 3.3 边界测试
func TestBadWordsCheckerEmptyWord(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("")
	c.Add("abc")

	assert.True(t, c.Contains("abc is here"))
	assert.False(t, c.Contains(""))
}

func TestBadWordsCheckerDuplicateAdd(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("abc")
	c.Add("abc")
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

	// "ab bc" 有空格边界
	result := c.Replace("ab bc", '*')
	assert.Equal(t, "** **", result)
}

func TestBadWordsCheckerACAutomaton(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"he", "she", "his", "hers"})

	// 英文：词边界检查防止误报
	assert.True(t, c.Contains("he said"))
	assert.True(t, c.Contains("she said"))
	assert.True(t, c.Contains("this is his book"))
	assert.True(t, c.Contains("hers"))
	assert.False(t, c.Contains("the"))  // "he" 在 "the" 内部
	assert.False(t, c.Contains("ushe")) // "she"/"he" 在 "ushe" 内部

	// 中文：AC 失败链接链正常工作
	c2 := NewBadWordsChecker()
	c2.AddWords([]string{"敏感词", "感词"})
	matches := c2.FindAll("这是敏感词文本")
	assert.Len(t, matches, 2)
}

func TestBadWordsCheckerBuildOnlyWhenDirty(t *testing.T) {
	c := NewBadWordsChecker()
	c.Add("test")

	assert.True(t, c.Contains("test sentence"))
	assert.True(t, c.Contains("another test"))
	assert.False(t, c.Contains("clean text"))
}

// 3.4 并发安全测试
func TestBadWordsCheckerConcurrentReadWrite(t *testing.T) {
	c := NewBadWordsChecker()
	c.AddWords([]string{"bad", "evil"})

	var wg sync.WaitGroup

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
			assert.True(t, c.Contains("a test case"))
		}()
	}
	wg.Wait()
}

// NormalizeCJK 测试
func TestNormalizeCJKRemoveSpaces(t *testing.T) {
	normalized, index := NormalizeCJK("违 禁 词")
	assert.Equal(t, "违禁词", normalized)
	assert.Equal(t, []int{0, 2, 4}, index)
}

func TestNormalizeCJKConsecutiveSpaces(t *testing.T) {
	normalized, index := NormalizeCJK("违   禁  词")
	assert.Equal(t, "违禁词", normalized)
	assert.Equal(t, []int{0, 4, 7}, index)
}

func TestNormalizeCJKPreserveASCIISpaces(t *testing.T) {
	normalized, index := NormalizeCJK("hello world")
	assert.Equal(t, "hello world", normalized)
	expected := make([]int, 11)
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, index)
}

func TestNormalizeCJKMixedText(t *testing.T) {
	normalized, _ := NormalizeCJK("hello 世 界")
	assert.Equal(t, "hello 世界", normalized)
}

func TestNormalizeCJKNoChange(t *testing.T) {
	normalized, index := NormalizeCJK("正常文本")
	assert.Equal(t, "正常文本", normalized)
	assert.Equal(t, []int{0, 1, 2, 3}, index)
}

func TestNormalizeCJKEmpty(t *testing.T) {
	normalized, index := NormalizeCJK("")
	assert.Equal(t, "", normalized)
	assert.Empty(t, index)
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
