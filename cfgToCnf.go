package main

import (
	"fmt"
	"os"
	"strconv"
)

// Converts a Context-Free Grammar into Chomsky Normal Form
func from_cfg_to_cnf(cfg *Grammar) *Grammar {
	fmt.Fprintln(os.Stdout, "Removing initial symbol: ", cfg.Initial)
	cnf := remove_initial(cfg)

	fmt.Fprintln(os.Stdout, "Binarizando producciones...")
	cnf = binarize_productions(cnf)

	fmt.Fprintln(os.Stdout, "Eliminar producciones epsilon...")
	cnf = delete_epsilon_productions(cnf)

	fmt.Fprintln(os.Stdout, "Eliminar producciones unitarias...")
	cnf = delete_unary_productions(cnf)

	fmt.Fprintln(os.Stdout, "Eliminando producciones/símbolos sin uso...")
	cnf = remove_useless_productions(cnf)

	return cnf
}

// Adds a new initial symbol to avoid issues with existing productions
func remove_initial(cfg *Grammar) *Grammar {
	newInitial := "S'"
	cfg.Productions[newInitial] = [][]string{{cfg.Initial}}
	cfg.Initial = newInitial
	return cfg
}

// Converts productions into binary form (Chomsky Normal Form)
func binarize_productions(cfg *Grammar) *Grammar {
	generator := construct_generator()
	for startState, transitions := range cfg.Productions {
		for transitionIdx, states := range transitions {
			cfg = bin_production(cfg, generator, states, startState, transitionIdx)
		}
	}
	return cfg
}

func bin_production(cfg *Grammar, generator func() string, transitionStates []string, startState string, transitionIdx int) *Grammar {
	if len(transitionStates) <= 1 {
		return cfg
	}

	if len(transitionStates) == 2 {
		allAreVariables := true
		for _, state := range transitionStates {
			if _, isTerminal := cfg.Terminals[state]; isTerminal {
				allAreVariables = false
			}
		}

		if allAreVariables {
			return cfg
		}
	}

	firstState := transitionStates[0]

	if _, firstIsTerminal := cfg.Terminals[firstState]; !firstIsTerminal {
		terminalSideState := add_new_state(generator, cfg, [][]string{{firstState}})
		nonTerminalSideState := add_new_state(generator, cfg, [][]string{transitionStates[1:]})

		cfg.Productions[startState][transitionIdx] = []string{terminalSideState, nonTerminalSideState}
		bin_production(cfg, generator, transitionStates[1:], nonTerminalSideState, 0)
	} else {
		stateOfOthers := add_new_state(generator, cfg, [][]string{transitionStates[1:]})
		cfg.Productions[startState][transitionIdx] = []string{firstState, stateOfOthers}
		bin_production(cfg, generator, transitionStates[1:], stateOfOthers, 0)
	}

	return cfg
}

// Adds a new state to the grammar and returns the newly created state
func add_new_state(generator func() string, cfg *Grammar, transitions [][]string) string {
	state := generator()
	_, transitionIsNotNew := cfg.Productions[state]
	for ; transitionIsNotNew; _, transitionIsNotNew = cfg.Productions[state] {
	}

	cfg.Productions[state] = transitions
	return state
}

// Generates a unique state name
func construct_generator() func() string {
	count := -1
	return func() string {
		count += 1
		return strconv.Itoa(count)
	}
}

// Deletes epsilon productions from the grammar
func delete_epsilon_productions(cfg *Grammar) *Grammar {
	statesToReplaceByEpsilonTransition := make(map[string][]string)

	for originalState, transitions := range cfg.Productions {
		for _, states := range transitions {
			for idx, stateValue := range states {
				if stateValue == "_" {
					firstPart := states[:idx]
					secondPart := states[idx+1:]
					firstPart = append(firstPart, secondPart...)
					statesToReplaceByEpsilonTransition[originalState] = firstPart
				}
			}
		}
	}

	for stateWithEpsilon, statesToAdd := range statesToReplaceByEpsilonTransition {
		for originalState, transitions := range cfg.Productions {
			for transitionIdx, states := range transitions {
				for stateWithEpsIdx, state := range states {
					if state == stateWithEpsilon {
						before := states[:stateWithEpsIdx]
						after := states[stateWithEpsIdx+1:]
						before = append(before, statesToAdd...)
						before = append(before, after...)

						cfg.Productions[originalState][transitionIdx] = before
					}
				}
			}
		}
	}

	return cfg
}

// Deletes unary productions from the grammar
func delete_unary_productions(cfg *Grammar) *Grammar {
	for originalState, transitions := range cfg.Productions {
		for _, states := range transitions {
			if len(states) != 1 {
				continue
			}

			unaryState := states[0]
			if _, isTerminal := cfg.Terminals[unaryState]; isTerminal {
				continue
			}

			// Asegúrate de que la eliminación de producciones unitarias no elimine las reglas principales
			relatedStates := cfg.Productions[unaryState]
			for _, transition := range relatedStates {
				if len(transition) == 1 && transition[0] == originalState {
					continue
				}
				cfg.Productions[originalState] = append(cfg.Productions[originalState], transition)
			}
		}
	}
	return cfg
}


// Deletes useless productions from the grammar
func remove_useless_productions(cfg *Grammar) *Grammar {
	usedStates := map[string]struct{}{
		"S'": {},
	}
	allStates := map[string]struct{}{}

	for originalState, transitions := range cfg.Productions {
		allStates[originalState] = struct{}{}

		for _, states := range transitions {
			for _, stateValue := range states {
				if _, isTerminal := cfg.Terminals[stateValue]; isTerminal {
					continue
				}

				usedStates[stateValue] = struct{}{}
			}
		}
	}

	uselessStates := map[string]struct{}{}
	for state := range allStates {
		if _, isUsed := usedStates[state]; !isUsed {
			uselessStates[state] = struct{}{}
		}
	}

	for uselessState := range uselessStates {
		delete(cfg.Productions, uselessState)
	}

	return cfg
}
