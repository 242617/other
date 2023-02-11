package leetcode

/*

([])()

*/

func IsValid(s string) bool {
	if len(s) <= 1 || len(s)%2 != 0 {
		return false
	}
	m := map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
	}
	if len(s) == 2 {
		e, ok := m[rune(s[0])]
		if ok && rune(s[1]) == e {
			return true
		}
	}
	stack := make([]rune, 0, len(s)/2)

	current := rune(s[0])
	closing, isOpening := m[current]
	if !isOpening {
		return false
	}
	stack = append(stack, closing)

	for i := 1; i < len(s); i++ {
		current = rune(s[i])
		closing, isOpening := m[current]
		if isOpening {
			stack = append(stack, closing)
		} else {
			if len(stack) == 0 {
				return false
			}
			check := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if current != check {
				return false
			}
		}
	}
	return len(stack) == 0
}
