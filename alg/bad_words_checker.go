package alg

import (
	"sort"
	"sync"
	"unicode/utf8"
)

// BadWordMatch banned word match result with position info
// BadWordMatch 违禁词匹配结果
type BadWordMatch struct {
	Word      string // 匹配到的违禁词原文
	ByteStart int    // 在原文中的字节起始位置
	ByteEnd   int    // 在原文中的字节结束位置（不包含）
}

// BadWordsChecker bad words checker based on Aho-Corasick automaton, O(n) time complexity multi-pattern matching
// BadWordsChecker 基于 AC 自动机的违禁词检查器，O(n) 时间复杂度多模式匹配
type BadWordsChecker struct {
	root       *TrieTreeNode
	lock       sync.RWMutex
	dirty      bool // 标记是否需要重建 AC 失败链接
	normalizer *BadWordsNormalizer
}

// NewBadWordsChecker create new bad words checker instance
// NewBadWordsChecker 创建违禁词检查器实例
func NewBadWordsChecker() *BadWordsChecker {
	return &BadWordsChecker{
		root:       &TrieTreeNode{},
		normalizer: &BadWordsNormalizer{},
	}
}

// Add add a single banned word, empty string will be ignored, duplicate word is idempotent
// Add 添加单个违禁词，空字符串忽略，重复添加幂等
func (c *BadWordsChecker) Add(word string) {
	if len(word) == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.addLocked(word)
}

// AddWords add multiple banned words at once
// AddWords 批量添加违禁词
func (c *BadWordsChecker) AddWords(words []string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, word := range words {
		if len(word) > 0 {
			c.addLocked(word)
		}
	}
}

// addLocked insert a word into trie tree, caller must hold write lock
// addLocked 将违禁词插入字典树，调用方需持有写锁
func (c *BadWordsChecker) addLocked(word string) {
	node := c.root
	for _, b := range []byte(word) {
		node = c.getOrCreateChild(node, b)
	}
	node.IsEnd = true
	c.dirty = true
}

// Build manually build AC automaton failure links, can be triggered automatically after Add
// Build 手动构建 AC 自动机失败链接，Add 后自动触发
func (c *BadWordsChecker) Build() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.dirty {
		return
	}
	c.build()
}

// Contains check if text contains any banned word, O(n) time complexity
// Contains 检查文本是否包含违禁词，O(n) 时间复杂度
func (c *BadWordsChecker) Contains(text string) bool {
	c.lock.RLock()
	empty := len(c.root.Children) == 0
	dirty := c.dirty
	c.lock.RUnlock()

	if empty {
		return false
	}
	if dirty {
		c.ensureBuilt()
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.scan(text, func(byteEnd int, node *TrieTreeNode) bool {
		return true // 首次匹配即停止
	})
}

// FindAll find all banned words in text with byte positions
// FindAll 查找文本中所有违禁词及字节位置
func (c *BadWordsChecker) FindAll(text string) []BadWordMatch {
	c.lock.RLock()
	empty := len(c.root.Children) == 0
	dirty := c.dirty
	c.lock.RUnlock()

	if empty {
		return nil
	}
	if dirty {
		c.ensureBuilt()
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	var matches []BadWordMatch
	c.scan(text, func(byteEnd int, node *TrieTreeNode) bool {
		byteStart := byteEnd - node.Depth
		word := text[byteStart:byteEnd]
		matches = append(matches, BadWordMatch{
			Word:      word,
			ByteStart: byteStart,
			ByteEnd:   byteEnd,
		})
		return false // 收集所有匹配
	})
	return matches
}

// Replace replace all banned words in text with mask character at rune level, overlapping matches are merged before replace
// Replace 将文本中违禁词替换为 mask 字符（rune 级别），重叠匹配会先合并
func (c *BadWordsChecker) Replace(text string, mask rune) string {
	matches := c.FindAll(text)
	if len(matches) == 0 {
		return text
	}

	// 合并重叠的匹配区间（贪心合并）
	merged := c.mergeByteRanges(matches)
	if len(merged) == 0 {
		return text
	}

	// 构建字节到 rune 位置的映射表
	runes := []rune(text)
	byteToRuneIndex := make([]int, len(text)+1)
	ri := 0
	for bi := 0; bi < len(text); {
		byteToRuneIndex[bi] = ri
		_, size := utf8.DecodeRuneInString(text[bi:])
		bi += size
		ri++
	}
	byteToRuneIndex[len(text)] = ri

	// 在 rune 级别执行替换
	for _, m := range merged {
		rStart := byteToRuneIndex[m.ByteStart]
		rEnd := byteToRuneIndex[m.ByteEnd]
		for i := rStart; i < rEnd; i++ {
			runes[i] = mask
		}
	}

	return string(runes)
}

// ensureBuilt ensure AC automaton is built, using double-checked locking
// ensureBuilt 确保 AC 自动机已构建，使用双重检查锁模式
func (c *BadWordsChecker) ensureBuilt() {
	c.lock.RLock()
	if !c.dirty {
		c.lock.RUnlock()
		return
	}
	c.lock.RUnlock()

	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.dirty {
		return
	}
	c.build()
}

// build BFS construction of AC automaton failure links, caller must hold write lock
// build BFS 构建 AC 自动机失败链接，调用方需持有写锁
func (c *BadWordsChecker) build() {
	// 确保根节点的基础状态
	c.root.Fail = nil

	// BFS 队列，从第 1 层开始
	queue := make([]*TrieTreeNode, 0, len(c.root.Children))
	for _, child := range c.root.Children {
		child.Fail = c.root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, child := range node.Children {
			// 沿父节点的失败链接查找匹配
			failNode := node.Fail
			for failNode != c.root && c.findChild(failNode, child.Data) == nil {
				failNode = failNode.Fail
			}
			if matched := c.findChild(failNode, child.Data); matched != nil && matched != child {
				child.Fail = matched
			} else {
				child.Fail = c.root
			}
			queue = append(queue, child)
		}
	}

	c.dirty = false
}

// scan AC automaton text scan, onMatch called for each matched word, return true to stop early
// scan AC 自动机扫描文本，onMatch 每匹配一个违禁词调用一次，返回 true 表示提前终止
func (c *BadWordsChecker) scan(text string, onMatch func(byteEnd int, node *TrieTreeNode) bool) bool {
	node := c.root

	for i := 0; i < len(text); i++ {
		b := text[i]

		// 回溯失败链接直到匹配或回到根节点
		for node != c.root && c.findChild(node, b) == nil {
			node = node.Fail
		}

		if child := c.findChild(node, b); child != nil {
			node = child
		}

		// 检查当前节点及其失败链接链上的所有匹配
		for temp := node; temp != c.root; temp = temp.Fail {
			if temp.IsEnd {
				byteEnd := i + 1
				byteStart := byteEnd - temp.Depth

				// ASCII-only 违禁词需通过词边界检查
				if !c.isWordBoundary(text, byteStart, byteEnd) {
					continue
				}

				if onMatch(byteEnd, temp) {
					return true
				}
			}
		}
	}

	return false
}

// isWordBoundary check if the match at [byteStart, byteEnd) in text is at a word boundary
// isWordBoundary 检查匹配位置是否在词边界，ASCII 词前后不能是 ASCII 字母
func (c *BadWordsChecker) isWordBoundary(text string, byteStart, byteEnd int) bool {
	word := text[byteStart:byteEnd]
	if !c.normalizer.isASCIIWord(word) {
		return true // 非 ASCII 违禁词保持子串匹配
	}
	if byteStart > 0 && c.normalizer.isASCIILetter(text[byteStart-1]) {
		return false
	}
	if byteEnd < len(text) && c.normalizer.isASCIILetter(text[byteEnd]) {
		return false
	}
	return true
}

// findChild binary search child node by data byte in sorted children slice
// findChild 在有序子节点切片中二分查找指定字节的子节点
func (c *BadWordsChecker) findChild(node *TrieTreeNode, data byte) *TrieTreeNode {
	if len(node.Children) == 0 {
		return nil
	}
	idx := sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Data >= data
	})
	if idx < len(node.Children) && node.Children[idx].Data == data {
		return node.Children[idx]
	}
	return nil
}

// getOrCreateChild find or create child node, new node's Fail defaults to root
// getOrCreateChild 查找或创建子节点，新节点的 Fail 默认为 root
func (c *BadWordsChecker) getOrCreateChild(node *TrieTreeNode, data byte) *TrieTreeNode {
	childCount := len(node.Children)
	if childCount > 0 {
		idx := sort.Search(childCount, func(i int) bool {
			return node.Children[i].Data >= data
		})
		if idx < childCount && node.Children[idx].Data == data {
			return node.Children[idx]
		}
	}
	// 创建新节点
	newNode := &TrieTreeNode{
		Data:  data,
		Fail:  c.root,
		Depth: node.Depth + 1,
	}
	node.Children = append(node.Children, newNode)
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Data < node.Children[j].Data
	})
	return newNode
}

// mergeByteRanges merge overlapping byte ranges by sorting and greedy merging
// mergeByteRanges 排序后贪心合并重叠的字节区间
func (c *BadWordsChecker) mergeByteRanges(matches []BadWordMatch) []BadWordMatch {
	if len(matches) <= 1 {
		return matches
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].ByteStart < matches[j].ByteStart
	})

	merged := []BadWordMatch{matches[0]}
	for i := 1; i < len(matches); i++ {
		last := &merged[len(merged)-1]
		if matches[i].ByteStart <= last.ByteEnd {
			if matches[i].ByteEnd > last.ByteEnd {
				last.ByteEnd = matches[i].ByteEnd
			}
		} else {
			merged = append(merged, matches[i])
		}
	}
	return merged
}

// NormalizeCJK remove ASCII spaces between CJK characters, delegation to BadWordsNormalizer
// NormalizeCJK 删除 CJK 字符之间的 ASCII 空格（委托给 BadWordsNormalizer）
func NormalizeCJK(text string) (normalized string, runeIndex []int) {
	return (&BadWordsNormalizer{}).NormalizeCJK(text)
}

// BadWordsNormalizer CJK/ASCII character judgment and text normalization utilities
// BadWordsNormalizer CJK/ASCII 字符判断与文本归一化工具
type BadWordsNormalizer struct{}

// isASCIILetter check if byte is an ASCII letter [a-zA-Z]
// isASCIILetter 检查字节是否为 ASCII 字母
func (n *BadWordsNormalizer) isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// isASCIIWord check if string is composed entirely of ASCII bytes
// isASCIIWord 检查字符串是否全为 ASCII 字节
func (n *BadWordsNormalizer) isASCIIWord(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 {
			return false
		}
	}
	return true
}

// isCJK check if rune is in CJK Unified Ideographs range
// isCJK 检查 rune 是否为 CJK 统一表意文字
func (n *BadWordsNormalizer) isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) ||
		(r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0x20000 && r <= 0x2A6DF)
}

// NormalizeCJK remove ASCII spaces between CJK characters to defeat space-based evasion
// NormalizeCJK 删除 CJK 字符之间的 ASCII 空格（用于反空格规避）
// 返回归一化文本和从归一化 rune 位置到原文 rune 位置的映射表
func (n *BadWordsNormalizer) NormalizeCJK(text string) (normalized string, runeIndex []int) {
	runes := []rune(text)

	if !n.hasCJKBoundSpace(runes) {
		runeIndex = make([]int, len(runes))
		for i := range runes {
			runeIndex[i] = i
		}
		return text, runeIndex
	}

	buf := make([]rune, 0, len(runes))
	runeIndex = make([]int, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' && n.isCJKBoundSpace(runes, i) {
			continue
		}
		buf = append(buf, runes[i])
		runeIndex = append(runeIndex, i)
	}

	return string(buf), runeIndex
}

// hasCJKBoundSpace returns true if any space in runes is between CJK characters
// hasCJKBoundSpace 检查是否存在被 CJK 字符包围的空格
func (n *BadWordsNormalizer) hasCJKBoundSpace(runes []rune) bool {
	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' && n.isCJKBoundSpace(runes, i) {
			return true
		}
	}
	return false
}

// isCJKBoundSpace returns true if runes[i] is a space bounded by CJK characters on both sides
// isCJKBoundSpace 检查位置 i 的空格是否前后有 CJK 字符包围（跳过连续空格）
func (n *BadWordsNormalizer) isCJKBoundSpace(runes []rune, i int) bool {
	if i <= 0 || i >= len(runes)-1 || runes[i] != ' ' {
		return false
	}
	// 向前查找最近的非空格字符
	hasCJKBefore := false
	for j := i - 1; j >= 0; j-- {
		if runes[j] != ' ' {
			hasCJKBefore = n.isCJK(runes[j])
			break
		}
	}
	if !hasCJKBefore {
		return false
	}
	// 向后查找最近的非空格字符
	for j := i + 1; j < len(runes); j++ {
		if runes[j] != ' ' {
			return n.isCJK(runes[j])
		}
	}
	return false
}
