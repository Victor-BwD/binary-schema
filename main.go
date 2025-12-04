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
	var bufferToSchema []byte

	for _, column := range s.Columns {
		columnBinary := column.Binary()
		bufferToSchema = append(bufferToSchema, columnBinary...)
	}

	return bufferToSchema
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
	buffer := make([]byte, 1024)
	offset := 0

	// Write column name
	nameBytesWritten := PutString(buffer[offset:], c.Name)
	offset += nameBytesWritten

	// Write column type
	typeBytesWritten, err := PutColumnType(buffer[offset:], c.Type, c.Nullable)
	if err != nil {
		fmt.Println("Error writing column type:", err)
		return nil
	}
	offset += typeBytesWritten

	return buffer[:offset]
}

func PutUvarint(slicedArray []byte, x uint64) []byte {
	if x == 0 {
		slicedArray[0] = 0xc0
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
	return slicedArray[:indice]
}

func PutString(buffer []byte, str string) int {
	length := len(str)

	fmt.Printf("DEBUG [PutString] Recebi a string: '%s'. Tamanho em bytes (len): %d\n", str, length)

	buffer[0] = 0xd1 // 0b11010001, indicando que é uma string

	lengthBytes := PutUvarint(buffer[1:], uint64(length)) // Passa a string e o tamanho dela para a função trazer o tamanho em bytes
	bytesWritten := len(lengthBytes)                      // Quantidade de bytes escritos para o tamanho da string

	copy(buffer[1+bytesWritten:], []byte(str)) // copia a string para o buffer, começando após o indice de bytes escritos

	fmt.Printf("DEBUG [PutString] Buffer preenchido até agora: %x\n", buffer[:1+bytesWritten+len(str)])
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
	schema := Schema{
		Columns: []Column{
			{Name: "id", Type: ColumnTypeInt, Nullable: false},
			{Name: "name", Type: ColumnTypeString, Nullable: false},
			{Name: "price", Type: ColumnTypeFloat, Nullable: true},
			{Name: "available", Type: ColumnTypeBool, Nullable: false},
		},
	}

	binario := schema.Binary() // Gera o binário do schema

	VisualizarBinario(binario)

	//fmt.Println(schema.JsonString())
}

func VisualizarBinario(data []byte) {
	fmt.Println("\n=== Mostra Bonito ===")

	for i := 0; i < len(data); i++ {
		b := data[i]

		if b == 0xd1 {
			fmt.Printf("\n[Nova Coluna] ")
		}

		if b == 0xd2 {
			fmt.Print("   |   [Definição] ")
		}
		fmt.Printf("%02x ", b)
	}
	fmt.Println("\n\n=========================")
}
