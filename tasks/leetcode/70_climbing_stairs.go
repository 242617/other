package leetcode

func ClimbStairs(n int) int {
	a := make([]int, 0, n+1)
	a = append(a, []int{0, 1, 2}...)
	for i := len(a); i <= n; i++ {
		a = append(a, a[i-1]+a[i-2])
	}
	return a[n]
}
