package search

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func main() {
	arr1 := make([]int, 100000)
	for i := 0; i < len(arr1); i++ {
		arr1[i] = rand.Intn(10000000)
	}
	fmt.Println(check(Max(arr1, 10)))

	arr2 := []int{0, 10, 100, 20, 1000, -30}
	fmt.Println(check(Max(arr2, 3)))

	arr3 := []int{}
	fmt.Println(check(Max(arr3, 2)))
}

func Max(arr []int, n int) []*int {
	item := NewItem(n)
	for _, i := range arr {
		item.check(i)
	}
	res := make([]*int, n)
	cur := item
	for i := 0; i < n; i++ {
		res[i] = cur.value
		cur = cur.child
	}
	return res
}

func check(arr []*int) string {
	strarr := make([]string, len(arr))
	for n, i := range arr {
		if i != nil {
			strarr[n] = strconv.Itoa(*i)
		} else {
			strarr[n] = "<nil>"
		}
	}
	return strings.Join(strarr, ", ")
}
