package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		RunCLI()
		return
	}

	RunAPI()
	return
}

func RunCLI() error {
	return nil
}

func RunAPI() error {
	return nil
}
