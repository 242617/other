package sort

import (
	"log"
	"reflect"
)

func Sort(array []interface{}, order Order) []interface{} {
	log.Println("Sort")

	result := make([]interface{}, len(array))
	for i, value := range array {
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)
		log.Println(v)

		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)

			tag := field.Tag.Get(order.FieldType)
			if tag == order.Field {
				log.Println(v.Field(j).Interface())
			}
		}

		result[i] = value
	}

	return result
}
