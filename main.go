package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"gopython/compiler"
	"gopython/runtime"
)

func main() {
	// flags
	bench := flag.Bool("bench", false, "benchmark program")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage:")
		fmt.Println("  go run . [--bench] [--cpuprofile file] test.py")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--- Compiling %s ---\n", filename)

	ast := compiler.Parse(string(src))
	chunk := compiler.Compile(ast)

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *bench {
		const iterations = 1_000_000

		fmt.Printf("--- Benchmark (%d iterations) ---\n", iterations)

		start := time.Now()

		vm := runtime.NewVM()

		for i := 0; i < iterations; i++ {
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
