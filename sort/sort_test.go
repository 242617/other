package sort

import (
	"log"
	"testing"
	"time"
)

type Item struct {
	Name      string    `json:"name"`
	Order     uint      `json:"order"`
	CreatedAt time.Time `json:"created_at"`
}
type Items []Item

var unsorted = Items{
	{"second", 14, time.Now().Add(3 * time.Second)},
	{"third", 20, time.Now().Add(time.Second)},
	{"first", 10, time.Now().Add(2 * time.Second)},
}

var tests = []struct {
	Order  Order
	Result []string
}{
	{
		Order: Order{
			Field:     "order",
			FieldType: "json",
			Direction: DirectionAscending,
		},
		Result: []string{"first", "second", "third"},
	},
	{
		Order: Order{
			Field:     "created_at",
			FieldType: "json",
			Direction: DirectionAscending,
		},
		Result: []string{"third", "first", "second"},
	},
}

func init() { log.SetFlags(log.Lshortfile) }

func TestSort(t *testing.T) {

	for _, test := range tests {

		data := make([]interface{}, len(unsorted))
		for i, item := range unsorted {
			data[i] = interface{}(item)
		}
		data = Sort(data, test.Order)
		sorted := make(Items, len(data))
		for i, value := range data {
			sorted[i] = value.(Item)
		}

		if len(sorted) != len(unsorted) {
			t.Fatalf("incorrect length: need %d, got %d", len(unsorted), len(sorted))
		}
		if sorted[0].Name != test.Result[0] {
			t.Fatalf(`expected "%s", but got "%s"`, test.Result[0], sorted[0].Name)
		}
		if sorted[1].Name != test.Result[1] {
			t.Fatalf(`expected "%s", but got "%s"`, test.Result[1], sorted[1].Name)
		}
		if sorted[2].Name != test.Result[2] {
			t.Fatalf(`expected "%s", but got "%s"`, test.Result[2], sorted[2].Name)
		}

	}

}
