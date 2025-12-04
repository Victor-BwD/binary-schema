package main

import (
	"reflect"
	"testing"
)

func TestPutString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "string simples ASCII",
			input:    "hello",
			expected: []byte{0xd1, 0xc1, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, // marca + length(1 byte, valor 5) + "hello"
		},
		{
			name:     "string com caractere UTF-8 'à'",
			input:    "à",
			expected: []byte{0xd1, 0xc1, 0x02, 0xc3, 0xa0}, // marca + length(1 byte, valor 2) + UTF-8 "à"
		},
		{
			name:     "string vazia",
			input:    "",
			expected: []byte{0xd1, 0xc0}, // marca + length(0)
		},
		{
			name:     "string com um caractere ASCII",
			input:    "a",
			expected: []byte{0xd1, 0xc1, 0x01, 0x61}, // marca + length(1 byte, valor 1) + "a"
		},
		{
			name:     "string com emoji",
			input:    "😀",
			expected: []byte{0xd1, 0xc1, 0x04, 0xf0, 0x9f, 0x98, 0x80}, // marca + length(1 byte, valor 4) + UTF-8 "😀"
		},
		{
			name:     "string com caracteres especiais",
			input:    "café",
			expected: []byte{0xd1, 0xc1, 0x05, 0x63, 0x61, 0x66, 0xc3, 0xa9}, // marca + length(1 byte, valor 5) + UTF-8 "café"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, 100)
			bytesWritten := PutString(buffer, tt.input)

			expectedLength := len(tt.expected)
			if bytesWritten != expectedLength {
				t.Errorf("PutString() retornou %d bytes, esperado %d", bytesWritten, expectedLength)
			}

			result := buffer[:bytesWritten]
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PutString() = %v, esperado %v", result, tt.expected)
				t.Errorf("Em hexadecimal: resultado = %x, esperado = %x", result, tt.expected)
			}
		})
	}
}

func TestPutUvarint(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected []byte
	}{
		{
			name:     "valor zero",
			input:    0,
			expected: []byte{0xc0},
		},
		{
			name:     "valor pequeno",
			input:    5,
			expected: []byte{0xc1, 0x05},
		},
		{
			name:     "valor médio",
			input:    255,
			expected: []byte{0xc1, 0xff},
		},
		{
			name:     "valor grande",
			input:    65000,
			expected: []byte{0xc2, 0xfd, 0xe8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, 10)
			result := PutUvarint(buffer, tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PutUvarint() = %v, esperado %v", result, tt.expected)
				t.Errorf("Em hexadecimal: resultado = %x, esperado = %x", result, tt.expected)
			}
		})
	}
}

func TestColumnTypeBinary(t *testing.T) {
	tests := []struct {
		name       string
		columnType ColumnType
		nullable   bool
		expected   uint8
		shouldErr  bool
	}{
		{
			name:       "ColumnTypeString não nullable",
			columnType: ColumnTypeString,
			nullable:   false,
			expected:   0x81, // 0b10000001
		},
		{
			name:       "ColumnTypeString nullable",
			columnType: ColumnTypeString,
			nullable:   true,
			expected:   0xc1, // 0b11000001 (0x81 | 0x40)
		},
		{
			name:       "ColumnTypeInt não nullable",
			columnType: ColumnTypeInt,
			nullable:   false,
			expected:   0x82, // 0b10000010
		},
		{
			name:       "ColumnTypeInt nullable",
			columnType: ColumnTypeInt,
			nullable:   true,
			expected:   0xc2, // 0b11000010 (0x82 | 0x40)
		},
		{
			name:       "ColumnTypeFloat não nullable",
			columnType: ColumnTypeFloat,
			nullable:   false,
			expected:   0x83, // 0b10000011
		},
		{
			name:       "ColumnTypeFloat nullable",
			columnType: ColumnTypeFloat,
			nullable:   true,
			expected:   0xc3, // 0b11000011 (0x83 | 0x40)
		},
		{
			name:       "ColumnTypeBool não nullable",
			columnType: ColumnTypeBool,
			nullable:   false,
			expected:   0x84, // 0b10000100
		},
		{
			name:       "ColumnTypeBool nullable",
			columnType: ColumnTypeBool,
			nullable:   true,
			expected:   0xc4, // 0b11000100 (0x84 | 0x40)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.columnType.Binary(tt.nullable)

			if tt.shouldErr && err == nil {
				t.Error("Esperava um erro mas não recebeu nenhum")
			}

			if !tt.shouldErr && err != nil {
				t.Errorf("Não esperava erro mas recebeu: %v", err)
			}

			if result != tt.expected {
				t.Errorf("ColumnType.Binary() = 0x%02x, esperado 0x%02x", result, tt.expected)
			}
		})
	}
}

func TestPutColumnType(t *testing.T) {
	tests := []struct {
		name       string
		columnType ColumnType
		nullable   bool
		expected   []byte
	}{
		{
			name:       "Int não nullable",
			columnType: ColumnTypeInt,
			nullable:   false,
			expected:   []byte{0xd2, 0x82},
		},
		{
			name:       "String nullable",
			columnType: ColumnTypeString,
			nullable:   true,
			expected:   []byte{0xd2, 0xc1},
		},
		{
			name:       "Float nullable",
			columnType: ColumnTypeFloat,
			nullable:   true,
			expected:   []byte{0xd2, 0xc3},
		},
		{
			name:       "Bool não nullable",
			columnType: ColumnTypeBool,
			nullable:   false,
			expected:   []byte{0xd2, 0x84},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, 10)
			bytesWritten, err := PutColumnType(buffer, tt.columnType, tt.nullable)

			if err != nil {
				t.Errorf("PutColumnType() erro = %v", err)
				return
			}

			if bytesWritten != 2 {
				t.Errorf("PutColumnType() escreveu %d bytes, esperado 2", bytesWritten)
			}

			result := buffer[:bytesWritten]
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PutColumnType() = %x, esperado %x", result, tt.expected)
			}
		})
	}
}

func TestColumnBinary(t *testing.T) {
	tests := []struct {
		name     string
		column   Column
		expected []byte
	}{
		{
			name: "Coluna id (int, não nullable)",
			column: Column{
				Name:     "id",
				Type:     ColumnTypeInt,
				Nullable: false,
			},
			// 0xd1 (marca string) + 0xc1 0x02 (length=2) + "id" (0x69 0x64) + 0xd2 0x82 (type int não nullable)
			expected: []byte{0xd1, 0xc1, 0x02, 0x69, 0x64, 0xd2, 0x82},
		},
		{
			name: "Coluna name (string, não nullable)",
			column: Column{
				Name:     "name",
				Type:     ColumnTypeString,
				Nullable: false,
			},
			// 0xd1 + 0xc1 0x04 (length=4) + "name" (0x6e 0x61 0x6d 0x65) + 0xd2 0x81
			expected: []byte{0xd1, 0xc1, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0xd2, 0x81},
		},
		{
			name: "Coluna price (float, nullable)",
			column: Column{
				Name:     "price",
				Type:     ColumnTypeFloat,
				Nullable: true,
			},
			// 0xd1 + 0xc1 0x05 (length=5) + "price" (0x70 0x72 0x69 0x63 0x65) + 0xd2 0xc3 (float nullable)
			expected: []byte{0xd1, 0xc1, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0xd2, 0xc3},
		},
		{
			name: "Coluna available (bool, não nullable)",
			column: Column{
				Name:     "available",
				Type:     ColumnTypeBool,
				Nullable: false,
			},
			// 0xd1 + 0xc1 0x09 (length=9) + "available" + 0xd2 0x84
			expected: []byte{0xd1, 0xc1, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0xd2, 0x84},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.column.Binary()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Column.Binary() falhou para '%s'", tt.name)
				t.Errorf("Resultado: %x", result)
				t.Errorf("Esperado:  %x", tt.expected)
			}
		})
	}
}

func TestSchemaBinary(t *testing.T) {
	schema := Schema{
		Columns: []Column{
			{Name: "id", Type: ColumnTypeInt, Nullable: false},
			{Name: "name", Type: ColumnTypeString, Nullable: false},
			{Name: "price", Type: ColumnTypeFloat, Nullable: true},
			{Name: "available", Type: ColumnTypeBool, Nullable: false},
		},
	}

	result := schema.Binary()

	// Montando o esperado concatenando cada coluna
	expected := []byte{
		// Coluna "id"
		0xd1, 0xc1, 0x02, 0x69, 0x64, 0xd2, 0x82,
		// Coluna "name"
		0xd1, 0xc1, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0xd2, 0x81,
		// Coluna "price"
		0xd1, 0xc1, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0xd2, 0xc3,
		// Coluna "available"
		0xd1, 0xc1, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0xd2, 0x84,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Schema.Binary() produziu resultado incorreto")
		t.Errorf("Resultado: %x", result)
		t.Errorf("Esperado:  %x", expected)

		// Análise detalhada para debugging
		t.Log("\nAnálise detalhada:")
		t.Logf("Tamanho resultado: %d bytes", len(result))
		t.Logf("Tamanho esperado:  %d bytes", len(expected))

		minLen := len(result)
		if len(expected) < minLen {
			minLen = len(expected)
		}

		for i := 0; i < minLen; i++ {
			if result[i] != expected[i] {
				t.Logf("Diferença no byte %d: resultado=0x%02x, esperado=0x%02x", i, result[i], expected[i])
			}
		}
	}
}

func TestColumnTypeString(t *testing.T) {
	tests := []struct {
		columnType ColumnType
		expected   string
	}{
		{ColumnTypeString, "string"},
		{ColumnTypeInt, "int"},
		{ColumnTypeFloat, "float"},
		{ColumnTypeBool, "bool"},
		{ColumnType(99), "unknown"},
	}

	for _, tt := range tests {
		result := tt.columnType.String()
		if result != tt.expected {
			t.Errorf("ColumnType(%d).String() = %q, esperado %q", tt.columnType, result, tt.expected)
		}
	}
}

func TestSchemaJsonString(t *testing.T) {
	schema := Schema{
		Columns: []Column{
			{Name: "id", Type: ColumnTypeInt, Nullable: false},
			{Name: "name", Type: ColumnTypeString, Nullable: true},
		},
	}

	result := schema.JsonString()

	// Verifica se é um JSON válido e contém os campos esperados
	if result == "" {
		t.Error("JsonString() retornou string vazia")
	}

	// Verificações básicas de conteúdo
	if !contains(result, `"id"`) {
		t.Error("JSON não contém campo 'id'")
	}
	if !contains(result, `"name"`) {
		t.Error("JSON não contém campo 'name'")
	}
	if !contains(result, `"int"`) {
		t.Error("JSON não contém tipo 'int'")
	}
	if !contains(result, `"string"`) {
		t.Error("JSON não contém tipo 'string'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
