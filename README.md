# Bencode

[![Build & Test](https://github.com/extintor/bencode/actions/workflows/build.yml/badge.svg)](https://github.com/extintor/bencode/actions/workflows/build.yml)
[![golangci-lint](https://github.com/extintor/bencode/actions/workflows/lint.yml/badge.svg)](https://github.com/extintor/bencode/actions/workflows/lint.yml)

## Overview

This Bencode implementation in Go is distinct in its approach to encoding and decoding data. It utilizes Go struct tags, allowing developers to define objects for serialization and deserialization in a manner akin to JSON. This feature enhances the ease of use and integration with Go's native data structures, providing a seamless and efficient way to work with Bencode data.

## Features

- Struct Tag-Based Serialization: Encode and decode Bencode data using Go struct tags, offering a familiar and intuitive approach similar to JSON serialization.
- Support for Complex Types: Effortlessly handle complex data structures, including nested structs and slices, ensuring versatility in data representation.
- Streamlined Integration: Designed to align naturally with Go's standard practices, making it easy to incorporate into Go projects.
- Full Bencode Type Support: Efficiently manage all standard Bencode types like strings, integers, lists, and dictionaries.

## Getting Started
### Prerequisites

Go (Version 1.21 or higher)

### Installation

```bash
go get github.com/extintor/bencode
```

### Usage
#### Decoding

```Go
package main

import (
	"fmt"
	"github.com/extintor/bencode"
)

func main() {
	type Profile struct {
		Name    string   `bencode:"name"`
		Age     uint64   `bencode:"age"`
		Hobbies []string `bencode:"hobbies"`
	}

	type User struct {
		Username string  `bencode:"username"`
		Email    string  `bencode:"email"`
		Profile  Profile `bencode:"profile"`
	}

	input := []byte("d8:username8:johndoe5:email15:john@example.com7:profiled4:name9:John Doe3:agei30e7:hobbiesl7:guitar6:travel5:cookingeee")
	var user User

	err := bencode.Decode(input, &user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded User: %+v\n", user)
}
```

#### Encoding

```Go
package main

import (
	"fmt"
	"github.com/extintor/bencode"
)

func main() {
	type Profile struct {
		Name    string   `bencode:"name"`
		Age     uint64   `bencode:"age"`
		Hobbies []string `bencode:"hobbies"`
	}

	type User struct {
		Username string  `bencode:"username"`
		Email    string  `bencode:"email"`
		Profile  Profile `bencode:"profile"`
	}

	user := User{
		Username: "johndoe",
		Email:    "john@example.com",
		Profile: Profile{
			Name:    "John Doe",
			Age:     30,
			Hobbies: []string{"guitar", "travel", "cooking"},
		},
	}

	encodedData, err := bencode.Encode(&user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encoded Data: %s\n", encodedData)
}
```

## Contributing

Contributions are welcome!

## License

This project is licensed under the MIT License - see the LICENSE file for details.
