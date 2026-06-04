package alg

import (
	"sort"
	"sync"
	"unicode/utf8"
)

// BadWordsTextNormalizer defines text preprocessing strategy for matching.
//
// Designed as an interface to support future normalization strategies
// (case folding, full-width to half-width, Unicode NFKD, etc.). Only one
// implementation exists today (BadWordsCjkNormalizer). The interface exists as a
// strategic extension point: it separates policy (BadWordsChecker's
// matching logic) from mechanism (how text is normalized), allowing
// deployments to swap normalizers without modifying checker internals.
//
// BadWordsTextNormalizer 定义匹配前的文本预处理策略。
// 设计为接口以支持未来替换归一化策略。虽然当前仅有一种实现，但接口将
// 匹配策略与归一化机制解耦，允许部署时替换而无需修改检查器内部。
type BadWordsTextNormalizer interface {
	// Normalize preprocesses raw text for matching.
	// normalized is the transformed text; runeMapping[i] is the rune
	// position in rawText that corresponds to normalized rune at index i.
	//
	// Normalize 预处理原始文本。
	// normalized 是变换后的文本；runeMapping[i] 是归一化后第 i 个 rune
	// 在 rawText 中对应的原始 rune 位置。
	Normalize(rawText string) (normalized string, runeMapping []int)
}

// BadWordMatch banned word match result with position info
// BadWordMatch 违禁词匹配结果
type BadWordMatch struct {
	Word      string // 匹配到的违禁词原文
	ByteStart int    // 在原文中的字节起始位置
	ByteEnd   int    // 在原文中的字节结束位置（不包含）
}

// BadWordsCheckerOpts configures a BadWordsChecker.
// BadWordsCheckerOpts 检查器配置
type BadWordsCheckerOpts struct {
	// Normalizers is an ordered chain of text normalizers applied before matching.
	// Normalizers 归一化器链，按序应用于匹配前的文本。
	Normalizers []BadWordsTextNormalizer
}

// BadWordsChecker bad words checker based on Aho-Corasick automaton, O(n) time complexity multi-pattern matching
// BadWordsChecker 基于 AC 自动机的违禁词检查器，O(n) 时间复杂度多模式匹配
type BadWordsChecker struct {
	root        *TrieTreeNode
	lock        sync.RWMutex
	dirty       bool // 标记是否需要重建 AC 失败链接
	normalizers []BadWordsTextNormalizer
}

// NewBadWordsChecker create a bad words checker.
// When opts is nil or Normalizers is empty, defaults to BadWordsCjkNormalizer.
//
// NewBadWordsChecker 创建违禁词检查器。
// opts 为 nil 或 Normalizers 为空时，默认使用 BadWordsCjkNormalizer。
func NewBadWordsChecker(opts *BadWordsCheckerOpts) *BadWordsChecker {
	c := &BadWordsChecker{
		root: &TrieTreeNode{},
	}
	if opts != nil && opts.Normalizers != nil {
		c.normalizers = opts.Normalizers
	} else {
		c.normalizers = []BadWordsTextNormalizer{&BadWordsCjkNormalizer{}}
	}
	return c
}

// AddNormalizer appends a normalizer to the chain.
// AddNormalizer 向归一化链追加一个归一化器。
func (c *BadWordsChecker) AddNormalizer(normalizer BadWordsTextNormalizer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.normalizers = append(c.normalizers, normalizer)
}

// RemoveNormalizer removes a normalizer from the chain.
// RemoveNormalizer 从归一化链中移除一个归一化器。
func (c *BadWordsChecker) RemoveNormalizer(normalizer BadWordsTextNormalizer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i, n := range c.normalizers {
		if n == normalizer {
			c.normalizers = append(c.normalizers[:i], c.normalizers[i+1:]...)
			return
		}
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
	return c.scan(c.normalizeText(text), func(byteEnd int, node *TrieTreeNode) bool {
		return true
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

	normalizedText := c.normalizeText(text)
	var matches []BadWordMatch
	c.scan(normalizedText, func(byteEnd int, node *TrieTreeNode) bool {
		byteStart := byteEnd - node.Depth
		word := normalizedText[byteStart:byteEnd]
		matches = append(matches, BadWordMatch{
			Word:      word,
			ByteStart: byteStart,
			ByteEnd:   byteEnd,
		})
		return false
	})
	return matches
}

// Replace replace all banned words in text with mask character at rune level, overlapping matches are merged before replace
// Replace 将文本中违禁词替换为 mask 字符（rune 级别），重叠匹配会先合并
func (c *BadWordsChecker) Replace(text string, mask rune) string {
	normalizedText := c.normalizeText(text)
	matches := c.FindAll(text)
	if len(matches) == 0 {
		return text
	}

	merged := c.mergeByteRanges(matches)
	if len(merged) == 0 {
		return text
	}

	runes := []rune(normalizedText)
	byteToRuneIndex := make([]int, len(normalizedText)+1)
	ri := 0
	for bi := 0; bi < len(normalizedText); {
		byteToRuneIndex[bi] = ri
		_, size := utf8.DecodeRuneInString(normalizedText[bi:])
		bi += size
		ri++
	}
	byteToRuneIndex[len(normalizedText)] = ri

	for _, m := range merged {
		rStart := byteToRuneIndex[m.ByteStart]
		rEnd := byteToRuneIndex[m.ByteEnd]
		for i := rStart; i < rEnd; i++ {
			runes[i] = mask
		}
	}

	return string(runes)
}

// ─── 内部方法 ───────────────────────────────────────────────────────────

func (c *BadWordsChecker) addLocked(word string) {
	node := c.root
	for _, b := range []byte(word) {
		node = c.getOrCreateChild(node, b)
	}
	node.IsEnd = true
	c.dirty = true
}

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

func (c *BadWordsChecker) build() {
	c.root.Fail = nil

	queue := make([]*TrieTreeNode, 0, len(c.root.Children))
	for _, child := range c.root.Children {
		child.Fail = c.root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, child := range node.Children {
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

func (c *BadWordsChecker) scan(text string, onMatch func(byteEnd int, node *TrieTreeNode) bool) bool {
	node := c.root

	for i := 0; i < len(text); i++ {
		b := text[i]

		for node != c.root && c.findChild(node, b) == nil {
			node = node.Fail
		}

		if child := c.findChild(node, b); child != nil {
			node = child
		}

		for temp := node; temp != c.root; temp = temp.Fail {
			if temp.IsEnd {
				byteEnd := i + 1
				byteStart := byteEnd - temp.Depth

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

// isWordBoundary checks if the match position is at a word boundary.
// ASCII-only words cannot have ASCII letters immediately before or after.
// Non-ASCII words (e.g. CJK) always pass (substring matching).
//
// isWordBoundary 检查匹配位置是否在词边界。
// 全 ASCII 违禁词前后不能是 ASCII 字母；非 ASCII 违禁词始终通过（子串匹配）。
func (c *BadWordsChecker) isWordBoundary(text string, byteStart, byteEnd int) bool {
	for i := byteStart; i < byteEnd; i++ {
		if text[i] >= 128 {
			return true // 非 ASCII → 子串匹配
		}
	}
	if byteStart > 0 && ((text[byteStart-1] >= 'a' && text[byteStart-1] <= 'z') || (text[byteStart-1] >= 'A' && text[byteStart-1] <= 'Z')) {
		return false
	}
	if byteEnd < len(text) && ((text[byteEnd] >= 'a' && text[byteEnd] <= 'z') || (text[byteEnd] >= 'A' && text[byteEnd] <= 'Z')) {
		return false
	}
	return true
}

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

func (c *BadWordsChecker) normalizeText(rawText string) string {
	text := rawText
	for _, n := range c.normalizers {
		text, _ = n.Normalize(text)
	}
	return text
}

// ─── BadWordsCjkNormalizer ─────────────────────────────────────────────────────────

// BadWordsCjkNormalizer removes ASCII spaces between CJK characters.
// BadWordsCjkNormalizer 删除 CJK 字符间的 ASCII 空格。
type BadWordsCjkNormalizer struct{}

// Normalize removes ASCII spaces between CJK characters.
// Normalize 删除 CJK 字符之间的 ASCII 空格。
func (n *BadWordsCjkNormalizer) Normalize(rawText string) (normalized string, runeMapping []int) {
	runes := []rune(rawText)

	if !n.hasCJKBoundSpace(runes) {
		runeMapping = make([]int, len(runes))
		for i := range runes {
			runeMapping[i] = i
		}
		return rawText, runeMapping
	}

	buf := make([]rune, 0, len(runes))
	runeMapping = make([]int, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' && n.isCJKBoundSpace(runes, i) {
			continue
		}
		buf = append(buf, runes[i])
		runeMapping = append(runeMapping, i)
	}

	return string(buf), runeMapping
}

func (n *BadWordsCjkNormalizer) isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) ||
		(r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0x20000 && r <= 0x2A6DF)
}

func (n *BadWordsCjkNormalizer) hasCJKBoundSpace(runes []rune) bool {
	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' && n.isCJKBoundSpace(runes, i) {
			return true
		}
	}
	return false
}

func (n *BadWordsCjkNormalizer) isCJKBoundSpace(runes []rune, i int) bool {
	if i <= 0 || i >= len(runes)-1 || runes[i] != ' ' {
		return false
	}
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
	for j := i + 1; j < len(runes); j++ {
		if runes[j] != ' ' {
			return n.isCJK(runes[j])
		}
	}
	return false
}
