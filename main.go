package main

import (
	"fmt"
	"os"
	"time"

	"gopython/compiler"
	"gopython/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run . <file.py>")
		fmt.Println("  go run . --bench <file.py>")
		os.Exit(1)
	}

	bench := false
	var filename string

	if os.Args[1] == "--bench" {
		if len(os.Args) < 3 {
			fmt.Println("Missing filename")
			os.Exit(1)
		}
		bench = true
		filename = os.Args[2]
	} else {
		filename = os.Args[1]
	}

	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--- Compiling %s ---\n", filename)

	// Compile once
	ast := compiler.Parse(string(src))
	chunk := compiler.Compile(ast)

	if bench {
		const iterations = 1_000_000

		fmt.Printf("--- Benchmark (%d iterations) ---\n", iterations)

		start := time.Now()

		for i := 0; i < iterations; i++ {
			vm := runtime.NewVM()
			vm.PushFrame(chunk)
			vm.Run()
		}

		elapsed := time.Since(start)

		fmt.Println("Total:", elapsed)
		fmt.Printf("Per run: %v\n", elapsed/time.Duration(iterations))
		fmt.Printf("Runs/sec: %.0f\n", float64(iterations)/elapsed.Seconds())
		return
	}

	fmt.Println("--- Running VM ---")

	vm := runtime.NewVM()
	vm.PushFrame(chunk)
	vm.Run()

	fmt.Println("--- Execution Finished ---")
	vm.Dump()
}