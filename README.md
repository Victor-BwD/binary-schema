## Go Binary Schema Parser

A lightweight, high-performance Go library designed to serialize database schemas into a compact custom binary format.

This project aims to replace verbose formats (like JSON) with a highly optimized byte stream, reducing payload size and improving parsing speed for schema definitions.

### 🚀 Purpose

In distributed systems or database internals, sending schema definitions (column names, types, constraints) via JSON can be unnecessary heavy. This project solves that by:

- Compactness: Converting strings and enums into optimized byte sequences.

- Speed: Using direct byte manipulation instead of reflection-heavy serializers.

- Custom Protocol: Implementing a variable-length integer logic and bitwise flags for column properties.

### 🛠 Protocol Specification
The binary protocol uses specific Byte Markers to identify data types and structures.

Column types are encoded in a single byte. We use bitwise operations to store both the Type and the Nullable status in the same byte.

Base Types:

- String: 0x81

- Int: 0x82

- Float: 0x83

- Bool: 0x84

- Nullable Flag: To make a column nullable, we apply a bitwise OR with 0x40 (01000000).

Example: Float (0x83) + Nullable (0x40) = 0xC3

### 🚧 Status
[x] Schema Definition: Structs and Enums defined.

[x] Serializer (Marshal): Converting Schema to Binary (schema.Binary()).

[ ] Deserializer (Parser): Converting Binary back to Schema struct (Work in Progress).

[ ] JSON Export: Converting the parsed schema back to JSON.
