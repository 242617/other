package sort

import (
	"log"
	"testing"
	"time"
)

type Item struct {
	Name      string
	Order     uint      `xml:"order"`
	Count     int32     `custom_tag:"count"`
	Amount    float64   `bson:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
type Items []Item

var unsorted = Items{
	{"second", 14, -17, 120.00, time.Now().Add(3 * time.Second)},
	{"third", 20, -24, 20.40, time.Now().Add(time.Second)},
	{"first", 10, 2, -40.20, time.Now().Add(2 * time.Second)},
}

var tests = []struct {
	Description string
	Order       Order
	Result      []string
}{
	{
		Description: "ascending sort by uint test",
		Order: Order{
			Field:     "order",
			FieldType: "xml",
			Direction: DirectionAscending,
		},
		Result: []string{"first", "second", "third"},
	},
	{
		Description: "descending sort by int32 test",
		Order: Order{
			Field:     "count",
			FieldType: "custom_tag",
			Direction: DirectionDescending,
		},
		Result: []string{"first", "second", "third"},
	},
	{
		Description: "descending sort by float64 test",
		Order: Order{
			Field:     "amount",
			FieldType: "bson",
			Direction: DirectionAscending,
		},
		Result: []string{"first", "third", "second"},
	},
	{
		Description: "ascending sort by time.Time test",
		Order: Order{
			Field:     "created_at",
			FieldType: "json",
			Direction: DirectionAscending,
		},
		Result: []string{"third", "first", "second"},
	},
	{
		Description: "descending sort by time.Time test",
		Order: Order{
			Field:     "created_at",
			FieldType: "json",
			Direction: DirectionDescending,
		},
		Result: []string{"second", "first", "third"},
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
			t.Fatalf("error in %s: incorrect length: need %d, got %d", test.Description, len(unsorted), len(sorted))
		}

		for i := range sorted {
			if sorted[i].Name != test.Result[i] {
				t.Fatalf(`error in %s: expected "%s", but got "%s"`, test.Description, test.Result[i], sorted[i].Name)
			}
		}

	}

}
