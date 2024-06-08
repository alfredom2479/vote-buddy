package main

import "fmt"

func checkError(description string, err error) {
	if err != nil {
		fmt.Println(description)
	}
}
