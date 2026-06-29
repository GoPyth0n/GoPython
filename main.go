// main.go
package main

import (
	"fmt"
	"os"

	"gopython/compiler"
	"gopython/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file.py>")
		os.Exit(1)
	}
	filename := os.Args[1]

	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--- Compiling %s ---\n", filename)
	
	// 1. Lex & Parse
	ast := compiler.Parse(string(src))
	
	// 2. Codegen
	chunk := compiler.Compile(ast)

	// 3. VM Execution
	fmt.Println("--- Running VM ---")
	vm := runtime.NewVM()
	
	// The new VM architecture requires pushing the initial chunk as the base frame
	vm.PushFrame(chunk)
	
	// Execute the frame stack
	vm.Run()
	
	fmt.Println("--- Execution Finished ---")
	
	// Utilize your new Dump method to inspect globals, the stack, and the return value
	vm.Dump()
}