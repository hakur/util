package alg

// MaxNonRepeatSubString max non repeat substring, based on slide window
// MaxNonRepeatSubString 最大不重复子串，基于滑动窗口实现
func MaxNonRepeatSubString(s string) (ss string) {
	var left int      // 当前左指针
	var lastLeft int  // 和max length相结合的最近一次左指针
	var right int     // 当前右指针
	var length int    // 当前字符串最大长度
	var maxLength int // 查找过程中的最大长度

	for i := range s {
		if s[left] == s[i] {
			if maxLength < length {
				lastLeft = left
				left = i
				maxLength = length
			}
		}

		right = i + 1
		length = right - left
	}

	return s[lastLeft : lastLeft+maxLength]
}
