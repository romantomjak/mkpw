package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go.1password.io/spg"
)

const (
	usage = `
Usage: mkpw [options]

  Generates a memorable, but secure password.

Options:

  -size=<n>
	Generate a password with <n> words (default: 5)

  -separator=<class>
	Separate components using <class> seperator. Can be one of: hyphen, space,
	comma, period, underscore, digit, none. (default: hyphen)

  -capitalize=<scheme>
	Capitalize password according to <scheme>. Scheme can be one of: none,
	first, all, random, one. (default: one)

  -entropy
    Show the entropy of the password recipe
`
)

func main() {
	var size int
	var separator, capitalize string
	var entropy bool

	flags := flag.NewFlagSet("mkpw", flag.ContinueOnError)
	flags.IntVar(&size, "size", 5, "size")
	flags.StringVar(&separator, "separator", "hyphen", "separator")
	flags.StringVar(&capitalize, "capitalize", "one", "capitalize")
	flags.BoolVar(&entropy, "entropy", false, "entropy")
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(usage))
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	sep, err := parseSeparator(separator)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cap, err := parseCapitalize(capitalize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	wl, err := spg.NewWordList(spg.AgileWords)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up word list: %s\n", err)
		os.Exit(1)
	}

	recipe := spg.NewWLRecipe(size, wl)
	recipe.Capitalize = cap
	recipe.SeparatorFunc = sep

	p, err := recipe.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating password: %s\n", err)
		os.Exit(1)
	}

	if entropy {
		fmt.Printf("Password: %s\nEntropy: %.3f\n", p, p.Entropy)
	} else {
		fmt.Printf("%s\n", p)
	}
}

func parseSeparator(separator string) (spg.SFFunction, error) {
	var f spg.SFFunction
	var err error

	switch separator {
	case "hyphen":
		f = makeSeparatorFunc("-")
		err = nil
		break
	case "space":
		f = makeSeparatorFunc(" ")
		err = nil
		break
	case "comma":
		f = makeSeparatorFunc(",")
		err = nil
		break
	case "dot":
		f = makeSeparatorFunc(".")
		err = nil
		break
	case "underscore":
		f = makeSeparatorFunc("_")
		err = nil
		break
	case "digit":
		f = spg.SFDigits1
		err = nil
		break
	case "none":
		f = spg.SFNone
		err = nil
		break
	default:
		f = nil
		err = fmt.Errorf("invalid separator: %q", separator)
	}

	return f, err
}

func makeSeparatorFunc(value string) spg.SFFunction {
	return func() (string, spg.FloatE) {
		return value, 0
	}
}

func parseCapitalize(capitalize string) (spg.CapScheme, error) {
	var scheme spg.CapScheme
	var err error

	switch capitalize {
	case "none":
		fallthrough
	case "first":
		fallthrough
	case "one":
		fallthrough
	case "all":
		fallthrough
	case "capitalize":
		fallthrough
	case "random":
		scheme = spg.CapScheme(capitalize)
		err = nil
		break
	default:
		scheme = ""
		err = fmt.Errorf("invalid capitalization scheme: %q", capitalize)
	}

	return scheme, err
}
