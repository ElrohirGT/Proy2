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
	// S  -> NP NV | NV
	//
	// NP -> 5 NV
	//
	// NV -> P
	//
	// P  -> 2
	//
	// Then we can check the initial state to start parsing.
	Productions map[string][][]string
	Initial     string
	Terminals   map[string]struct{}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	rules := make(map[string][][]string)

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
			rules[first] = [][]string{}
		}

		transitions := strings.Split(rest, "|")
		for _, transition := range transitions {
			trimmed := strings.TrimSpace(transition)
			states := strings.Split(trimmed, " ")
			rules[first] = append(rules[first], states)
			fmt.Fprintln(os.Stdout, "Adding rule:", first, "->", trimmed)
		}
	}

	terminals := map[string]struct{}{}
	for _, transitions := range rules {
		for _, states := range transitions {
			for _, state := range states {
				if _, notTerminal := rules[state]; notTerminal {
					continue
				}

				terminals[state] = struct{}{}

			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	grammar := Grammar{
		Productions: rules,
		Terminals:   terminals,
		Initial:     "S",
	}
	fmt.Printf("%v\n", grammar)

	chomsky := from_cfg_to_cnf(&grammar)
	fmt.Printf("%v\n", chomsky)

}
