package sort

import (
	"reflect"
	"sort"
	"time"
)

func Sort(array []interface{}, order Order) []interface{} {
	c := container{array, order}
	sort.Sort(c)
	return c.d
}

type container struct {
	d []interface{}
	o Order
}

func (c container) Len() int      { return len(c.d) }
func (c container) Swap(i, j int) { c.d[i], c.d[j] = c.d[j], c.d[i] }
func (c container) Less(i, j int) bool {
	ti, vi := c.get(i)
	tj, vj := c.get(j)

	if ti != tj {
		panic("incompatible type")
	}

	var result bool
	switch vi.(type) {
	case uint:
		result = vi.(uint) < vj.(uint)
	case uint32:
		result = vi.(uint32) < vj.(uint32)
	case uint64:
		result = vi.(uint64) < vj.(uint64)
	case int:
		result = vi.(int) < vj.(int)
	case int32:
		result = vi.(int32) < vj.(int32)
	case int64:
		result = vi.(int64) < vj.(int64)
	case float32:
		result = vi.(float32) < vj.(float32)
	case float64:
		result = vi.(float64) < vj.(float64)
	case time.Time:
		timei, timej := vi.(time.Time), vj.(time.Time)
		result = timei.Before(timej)
	default:
		panic("unsopported type")
	}

	if c.o.Direction != DirectionAscending {
		result = !result
	}

	return result
}

func (c container) get(i int) (reflect.Type, interface{}) {
	t := reflect.TypeOf(c.d[i])
	v := reflect.ValueOf(c.d[i])
	for j := 0; j < t.NumField(); j++ {
		field := t.Field(j)
		tag := field.Tag.Get(c.o.FieldType)
		if tag == c.o.Field {
			return field.Type, v.Field(j).Interface()
		}
	}
	return nil, nil
}
