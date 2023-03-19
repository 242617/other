package leetcode

func Subsets(nums []int) [][]int {
	var pow func(a, p int) int
	pow = func(a, p int) int {
		if a == 0 {
			return 0
		}
		if p == 0 {
			return 1
		}
		if p == 1 {
			return a
		}
		return a * pow(a, p-1)
	}

	res := make([][]int, pow(2, len(nums)))
	for i := 0; i < pow(2, len(nums)); i++ {
		if res[i] == nil {
			res[i] = make([]int, 0, len(nums))
		}
		for j := 0; j < len(nums); j++ {
			if i&(1<<j) != 0 {
				res[i] = append(res[i], nums[j])
			}
		}
	}

	return res
}
