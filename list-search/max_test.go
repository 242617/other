package main

import "testing"

func TestMax(t *testing.T) {
	arr := Max([]int{0, 10, 100, 20, 1000, -30}, 3)
	if (arr[0] == nil && *arr[0] != 1000) &&
		(arr[1] == nil && *arr[1] != 100) &&
		(arr[2] == nil && *arr[2] != 20) {
		t.Fail()
	}
}
