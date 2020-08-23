package main

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func main() {
	// Creating UUID Version 4
	// panic on error
  var err error
	u1 := uuid.Must(uuid.NewV4(), err)
	fmt.Printf("UUIDv4: %s\n", u1)

	// or error handling
	u2 := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("UUIDv4: %s\n", u2)

	// Parsing UUID from string input
	u3, err := uuid.FromString("38400000-8cf0-11bd-b23e-10b96e4ef00d")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u3)
}
