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

func PutUvarint(slicedArray []byte, x uint64) []byte {
	if x == 0 {
		slicedArray[0] = 0xc0
		fmt.Printf("0x%02x ", slicedArray[0])
		return slicedArray[:1]
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
		fmt.Printf("0x%02x ", slicedArray[i])
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

	return slicedArray[:indice]
}

func PutString(buffer []byte, str string) int {
	length := byte(len(str))

	fmt.Printf("Tamanho da string: %d\n", byte(length))

	buffer[0] = 0xd1 // 0b11010001, indicando que é uma string
	fmt.Printf("0x%02x ", buffer[0])

	lengthBytes := PutUvarint(buffer[1:], uint64(length)) // Passa a string e o tamanho dela para a função trazer o tamanho em bytes
	bytesWritten := len(lengthBytes)

	copy(buffer[1+bytesWritten:], []byte(str)) // copia a string para o buffer, começando após o indice de bytes escritos

	strStart := 1 + bytesWritten // Calcula o índice onde a string começa no buffer
	for i := strStart; i < strStart+len(str); i++ {
		fmt.Printf("0x%02x ", buffer[i])
	}

	return 1 + bytesWritten + len(str)
}

func PutColumnType(buffer []byte, columnType ColumnType, nullable bool) (int, error) {
	value, err := columnType.Binary(nullable)
	if err != nil {
		return 0, err
	}

	buffer[0] = 0xd2
	buffer[1] = value
	fmt.Printf("0x%02x 0x%02x ", buffer[0], buffer[1])

	return 2, nil
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

	//PutUvarint(make([]byte, 10), 65000)
	array := make([]byte, 100)

	bytesInString := PutString(array, "à")

	fmt.Printf("\nBytes usados: %d\n", bytesInString)

	//fmt.Println(schema.JsonString())
}
