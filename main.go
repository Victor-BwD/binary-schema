package main

import (
	"encoding/json"
	"fmt"
)

type ColumnType int

const (
	ColumnTypeString ColumnType = iota
	ColumnTypeInt
	ColumnTypeFloat
	ColumnTypeBool
)

type Column struct {
	Name     string     `json:"name"`
	Type     ColumnType `json:"type"`
	Nullable bool       `json:"nullable"`
}

type Schema struct {
	Columns []Column `json:"columns"`
}

func (c ColumnType) String() string {
	switch c {
	case ColumnTypeString:
		return "string"
	case ColumnTypeInt:
		return "int"
	case ColumnTypeFloat:
		return "float"
	case ColumnTypeBool:
		return "bool"
	default:
		return "unknown"
	}
}

func (c ColumnType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (s *Schema) JsonString() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}

func (s *Schema) Binary() []uint8 {
	return make([]uint8, 0)
}

func (c ColumnType) Binary(nullable bool) (uint8, error) {
	var value uint8

	switch c {
	case ColumnTypeString:
		// value = 0b10000001
		value = 0x81 // 0b10000001
	case ColumnTypeInt:
		// value = 0b10000010
		value = 0x82 // 0b10000010
	case ColumnTypeFloat:
		// value = 0b10000011
		value = 0x83 // 0b10000011
	case ColumnTypeBool:
		value = 0x84 // 0b10000100
	default:
		return 0, fmt.Errorf("unknown column type: %v", c)
	}

	if nullable {
		value = value | 0x40 // Set the nullable bit (0b01000000)
	}

	return value, nil
}

func (c *Column) Binary() []uint8 {
	return make([]uint8, 0)
}

func PutUvarint(slicedArray []byte, x uint64) int {

	var size int // Tamanho do número em bytes
	if x == 0 {
		size = 1
	} else {
		for v := x; v > 0; v >>= 8 {
			size++
		}
	}

	fmt.Printf("Size: %d\n", size)
	return 0
}

func main() {
	/*schema := Schema{
		Columns: []Column{
			{Name: "id", Type: ColumnTypeInt, Nullable: false},
			{Name: "name", Type: ColumnTypeString, Nullable: false},
			{Name: "price", Type: ColumnTypeFloat, Nullable: true},
			{Name: "available", Type: ColumnTypeBool, Nullable: false},
		},
	}*/

	PutUvarint(make([]byte, 10), 66000)

	//fmt.Println(schema.JsonString())
}
