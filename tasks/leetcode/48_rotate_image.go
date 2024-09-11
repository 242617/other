package leetcode

/*
1 2 -> 3 1
3 4    4 2

1 2 3    7 4 1
4 5 6 -> 8 5 2
7 8 9    9 6 3

01 02 03 04      13 09 05 01
05 06 07 08  ->  14 10 06 02
09 10 11 12      15 11 07 03
13 14 15 16      16 12 08 04

01 02 03 04 05      21 16 11 06 01
06 07 08 09 10      22 17 12 07 02
11 12 13 14 15  ->  23 18 13 08 03
16 17 18 19 20      24 19 14 09 04
21 22 23 24 25      25 20 15 10 05
*/

func Rotate(matrix [][]int) {
	length := len(matrix)
	levels := length/2 + length%2
	center := levels - 1

	type point struct {
		x, y int
		val  int
	}

	var level int
	if length%2 == 1 {
		level++ // Skip center
	}

	for ; level < levels; level++ {
		start := center - level
		width := (2*level + 1) + (length+1)%2

		circle := make([]point, 0, 2*width+2*(width-2))
		x, y := start, start
		for ; x < start+width-1; x++ {
			circle = append(circle, point{x: x, y: y, val: matrix[y][x]})
		}
		for ; y < start+width-1; y++ {
			circle = append(circle, point{x: x, y: y, val: matrix[y][x]})
		}
		for ; x > start; x-- {
			circle = append(circle, point{x: x, y: y, val: matrix[y][x]})
		}
		for ; y > start; y-- {
			circle = append(circle, point{x: x, y: y, val: matrix[y][x]})
		}

		for i := 0; i < len(circle); i++ {
			x, y := circle[i].x, circle[i].y
			shift := i - (width - 1) // Clockwise
			if shift < 0 {
				shift += len(circle)
			}
			// if shift > len(circle)-1 { // Counter clockwise
			// 	shift %= len(circle)
			// }
			matrix[y][x] = circle[shift].val
		}
	}
}
