package main

import (
	"fmt"
	"os"
)

// Converts a Context Free Grammar into a Chomsky Normal Form
func from_cfg_to_cnf(cfg *Grammar) *Grammar {
	fmt.Fprintln(os.Stdout, "Removing initial simbol: ", cfg.Initial)
	cnf := remove_initial(cfg)
	cnf = binarize_productions(cnf)
	cnf = delete_epsilon(cnf)
	cnf = delete_unit(cnf)
	cnf = delete_useless(cnf)

	return cnf
}

func remove_initial(cfg *Grammar) *Grammar {
	// TODO: Implement
}

func binarize_productions(cfg *Grammar) *Grammar {
	// TODO: Implement
}

func delete_epsilon(cfg *Grammar) *Grammar {
	// TODO: Implement
}

func delete_unit(cfg *Grammar) *Grammar {
	// TODO: Implement
}

func delete_useless(cfg *Grammar) *Grammar {
	// TODO: Implement
}
