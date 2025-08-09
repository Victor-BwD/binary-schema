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
