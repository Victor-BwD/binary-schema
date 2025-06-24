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

	var marker byte
	switch size {
	case 1:
		marker = 0xc1
	case 2:
		marker = 0xc2
	case 3, 4:
		marker = 0xc4
		size = 4
	case 5, 6, 7, 8:
		marker = 0xc8
		size = 8
	}

	slicedArray[0] = marker

	for i := 1; i <= size; i++ {
		shift := (size - i) * 8
		slicedArray[i] = byte(x >> shift)
	}

	totalBytes := size + 1
	for i := 0; i < totalBytes; i++ {
		fmt.Printf("0b%08b ", slicedArray[i])
		if i < totalBytes-1 {
			fmt.Print(", ")
		}
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

	PutUvarint(make([]byte, 10), 500)

	//fmt.Println(schema.JsonString())
}
