package leetcode

func Generate(numRows int) [][]int {
	res := make([][]int, numRows)
	for i := 0; i < numRows; i++ {
		res[i] = make([]int, i+1)
		for j := 0; j < i+1; j++ {
			res
		}
	}
	return res
}
