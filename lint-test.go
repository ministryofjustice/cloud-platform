//go:build ignore

package main

import "fmt"

func badlyFormatted() {
	x := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range x {
		fmt.Println(k, v)
	}
}
