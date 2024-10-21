package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Grammar struct {
	// The rules of the grammar. For example:
	//
	// S  -> NP | NV
	//
	// NP -> 5
	//
	// NV -> P
	//
	// P  -> 2
	//
	// Then we can check the initial state to start parsing.
	Rules   map[string][]string
	Initial string
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	rules := make(map[string][]string)

	fmt.Println("Hello! Please input your CFG rules:")

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			fmt.Println("Input ended!")
			break
		}

		first, rest, separatorFound := strings.Cut(line, "->")

		if !separatorFound {
			fmt.Fprintln(os.Stdout, "Rule", line, "has invalid format! Remember to add ->")
			continue
		}

		first = strings.TrimSpace(first)
		rest = strings.TrimSpace(rest)
		_, hasKey := rules[first]

		if !hasKey {
			rules[first] = []string{}
		}

		states := strings.Split(rest, "|")
		for _, state := range states {
			trimmed := strings.TrimSpace(state)
			rules[first] = append(rules[first], trimmed)
			fmt.Fprintln(os.Stdout, "Adding rule:", first, "->", trimmed)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
