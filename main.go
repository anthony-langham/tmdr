package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	
	"github.com/anthonylangham/tmdr/internal/acronym"
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

	// Load the acronym repository
	repo, err := acronym.NewCSVRepository("data/acronyms.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading acronym database: %v\n", err)
		os.Exit(1)
	}

	if *randomFlag {
		a, err := repo.Random()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting random acronym: %v\n", err)
			os.Exit(1)
		}
		printAcronym(a)
		os.Exit(0)
	}

	// Look up the provided acronym (case-insensitive)
	acronymStr := strings.ToUpper(flag.Arg(0))
	a, err := repo.Find(acronymStr)
	if err != nil {
		fmt.Printf("Acronym '%s' not found.\n", acronymStr)
		fmt.Println("Try 'tmdr --help' for usage information.")
		os.Exit(1)
	}
	
	printAcronym(a)
}

func printAcronym(a *acronym.Acronym) {
	fmt.Printf("%s → %s\n", a.Acronym, a.FullForm)
	if a.Definition != "" {
		fmt.Println(a.Definition)
	}
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