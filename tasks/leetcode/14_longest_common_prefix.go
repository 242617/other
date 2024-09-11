package leetcode

import (
	"fmt"
	"log"
	"strings"
)

func LongestCommonPrefix(strs []string) string {
	// var res string
	// if len(strs) == 0 {
	// 	return res
	// }
	// if len(strs) == 1 {
	// 	return strs[0]
	// }
	// first := strs[0]
	// for i := 0; i < len(first); i++ {
	// 	for j := 1; j < len(strs); j++ {
	// 		if i >= len(strs[j]) || first[i] != strs[j][i] {
	// 			return res
	// 		}
	// 	}
	// 	res = res + string(first[i])
	// }
	// return res

	lists := make(map[*list]struct{}, len(strs))
	min := len(strs[0])
	for _, str := range strs {
		lists[stringToList(str)] = struct{}{}
		if len(str) < min {
			min = len(str)
		}
	}

	// for list := range lists {
	// 	log.Println(list)
	// }

	log.Println("min", min)
	for i := 0; i < min; i++ {
		log.Println(i)
		for list := range lists {
			// 		if list.length-1 < i {
			// 			delete(lists, list)
			// 			continue
			// 		}

			log.Println(list)

		}
	}
	return ""
}

func stringToList(str string) *list {
	if str == "" {
		return nil
	}
	var current *node
	for i := len(str) - 1; i >= 0; i-- {
		n := &node{
			next:  current,
			value: rune(str[i]),
		}
		current = n
	}
	return &list{
		head:   current,
		length: len(str),
	}
}

type list struct {
	head   *node
	length int
}

type node struct {
	next  *node
	value rune
}

func (l *list) String() string {
	var nodes []string
	for current := l.head; current != nil; current = current.next {
		nodes = append(nodes, string(current.value))
	}
	return fmt.Sprintf("list: { length: %d, nodes: %s }", l.length, strings.Join(nodes, "->"))
}
