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

// Converts a Context Free Grammar into a Chomsky Normal Form
func from_cfg_to_cnf() {
	// let the input be a string I consisting of n characters: a1 ... an.
	// let the grammar contain r nonterminal symbols R1 ... Rr, with start symbol R1.
	// let P[n,n,r] be an array of booleans. Initialize all elements of P to false.
	// let back[n,n,r] be an array of lists of backpointing triples. Initialize all elements of back to the empty list.
	//
	// for each s = 1 to n
	//
	//	for each unit production Rv → as
	//	    set P[1,s,v] = true
	//
	// for each l = 2 to n -- Length of span
	//
	//	for each s = 1 to n-l+1 -- Start of span
	//	    for each p = 1 to l-1 -- Partition of span
	//	        for each production Ra    → Rb Rc
	//	            if P[p,s,b] and P[l-p,s+p,c] then
	//	                set P[l,s,a] = true,
	//	                append <p,b,c> to back[l,s,a]
	//
	// if P[n,1,1] is true then
	//
	//	I is member of language
	//	return back -- by retracing the steps through back, one can easily construct all possible parse trees of the string.
	//
	// else
	//
	//	return "not a member of language"
}
