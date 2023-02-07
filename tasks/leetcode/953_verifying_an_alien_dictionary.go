package leetcode

import "sort"

func IsAlienSorted(words []string, order string) bool {
	alphabet := make(map[byte]int, len(order)) // Prepare alphabet
	for i := 0; i < len(order); i++ {
		alphabet[order[i]] = i
	}
	sorted := make([]string, len(words))
	copy(sorted, words)
	sort.SliceStable(sorted, func(i, j int) bool { // Sort words
		left, right := sorted[i], sorted[j]
		min := len(left)
		if len(right) < min {
			min = len(right)
		}
		for k := 0; k < min; k++ {
			if alphabet[left[k]] == alphabet[right[k]] {
				continue
			}
			return alphabet[left[k]] < alphabet[right[k]]
		}
		return len(left) == len(right) || len(left) < len(right)
	})
	for i := 0; i < len(words); i++ { // Compare sorted and initial slices
		if words[i] != sorted[i] {
			return false
		}
	}
	return true
}
