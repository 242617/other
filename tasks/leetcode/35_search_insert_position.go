package leetcode

func SearchInsert(nums []int, target int) int {
	if len(nums) == 0 {
		return 0
	}
	if target <= nums[0] {
		return 0
	}
	if target > nums[len(nums)-1] {
		return len(nums)
	}
	left, right := 0, len(nums)-1
	for current := (left + right) / 2; ; current = (left + right) / 2 {
		switch {
		case right-left == 1:
			switch {
			case target == nums[left]:
				return left

			case target == nums[right]:
				return right

			case target < nums[left]:
				if left == 0 {
					return left
				}
				return left - 1

			case target > nums[right]:
				return right + 1

			default: // target > nums[left] && target < nums[right]
				return left + 1

			}
		case left == right:
			return left

		case nums[current] > target:
			right = current
			continue

		case nums[current] < target:
			left = current
			continue

		case nums[current] == target:
			return current

		}
	}
}
