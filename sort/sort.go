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

	switch vi.(type) {
	case uint, uint32, uint64, int, int32, int64:
		uinti, uintj := vi.(uint), vj.(uint)
		return uinti < uintj
	case time.Time:
		timei, timej := vi.(time.Time), vj.(time.Time)
		return timei.Before(timej)
	}

	return false
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
