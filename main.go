package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const version = "0.1.0"

func main() {
	var (
		versionFlag = flag.Bool("version", false, "Show version information")
		helpFlag    = flag.Bool("help", false, "Show help information")
		randomFlag  = flag.Bool("random", false, "Display a random acronym")
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("tmdr version %s\n", version)
		os.Exit(0)
	}

	if *helpFlag || (flag.NArg() == 0 && !*randomFlag) {
		printHelp()
		os.Exit(0)
	}

	if *randomFlag {
		fmt.Println("Random mode: Coming soon!")
		os.Exit(0)
	}

	acronym := strings.ToUpper(flag.Arg(0))
	fmt.Printf("Looking up: %s (Coming soon!)\n", acronym)
}

func printHelp() {
	fmt.Println("tmdr - Too Medical; Didn't Read")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tmdr <acronym>     Look up a medical acronym")
	fmt.Println("  tmdr --random      Display a random acronym")
	fmt.Println("  tmdr --version     Show version information")
	fmt.Println("  tmdr --help        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tmdr abg           Look up ABG (Arterial Blood Gas)")
	fmt.Println("  tmdr hiv           Look up HIV (Human Immunodeficiency Virus)")
}