package sort

import (
	"log"
	"testing"
)

type Item struct {
	Name  string `json:"name"`
	Order uint   `json:"order"`
}
type Items []Item

var unsorted = Items{
	{"second", 14},
	{"third", 20},
	{"first", 10},
}
var order = Order{
	Field:     "order",
	FieldType: "json",
	Direction: DirectionAscending,
}

func init() { log.SetFlags(log.Lshortfile) }

func TestSort(t *testing.T) {

	data := make([]interface{}, len(unsorted))
	for i, item := range unsorted {
		data[i] = interface{}(item)
	}
	data = Sort(data, order)

	sorted := make(Items, len(data))
	for i, value := range data {
		sorted[i] = value.(Item)
	}

	if len(sorted) != len(unsorted) {
		t.Fatalf("incorrect length: need %d, got %d", len(unsorted), len(sorted))
	}
	if sorted[0].Name != "first" {
		t.Fatalf(`expected "%s", but got "%s"`, "first", sorted[0].Name)
	}
	if sorted[1].Name != "second" {
		t.Fatalf(`expected "%s", but got "%s"`, "second", sorted[1].Name)
	}
	if sorted[2].Name != "third" {
		t.Fatalf(`expected "%s", but got "%s"`, "third", sorted[2].Name)
	}

}
