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
	if x == 0 {
		slicedArray[0] = 0xc0
		fmt.Printf("0b%08b ", slicedArray[0])
		return 1
	}

	indice := 1

	for j := 56; j >= 0; j -= 8 {
		b := byte(x >> j)
		if b != 0 || j == 0 {
			slicedArray[indice] = b
			indice++
		}
	}

	size := indice - 1
	slicedArray[0] = 0xc0 | byte(size)
	for i := 0; i < indice; i++ {
		fmt.Printf("0b%08b ", slicedArray[i])
	}

	/*if x <= 0xff {
		slicedArray[0] = byte(x)
		fmt.Printf("0b%08b\n", slicedArray[0])
		return 1
	}
	if x <= 0xffff {
		slicedArray[0] = byte(x)
		slicedArray[1] = byte(x >> 8)
		fmt.Printf("0b%08b, 0b%08b\n", slicedArray[0], slicedArray[1])
		return 2
	}
	if x <= 0xffffffff {
		slicedArray := make([]byte, 5)

		slicedArray[0] = 0xc0 | 4
		slicedArray[1] = byte(x >> 24)
		slicedArray[2] = byte(x >> 16)
		slicedArray[3] = byte(x >> 8)
		slicedArray[4] = byte(x)
		fmt.Printf("0b%08b, 0b%08b, 0b%08b, 0b%08b, 0b%08b\n",
			slicedArray[0], slicedArray[1], slicedArray[2], slicedArray[3], slicedArray[4])
		return 4
	}*/

	return indice
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

	PutUvarint(make([]byte, 10), 65000)

	//fmt.Println(schema.JsonString())
}
