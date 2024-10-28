package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Grammar struct {
	Productions map[string][][]string
	Initial     string
	Terminals   map[string]struct{}
}

func main() {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("✨ Bienvenido al Verificador de Frases ✨")
	fmt.Println("════════════════════════════════════════")

	// Leer la gramática desde el archivo
	grammarFile := "input.txt"
	data := readFile(grammarFile)

	fmt.Println("📖  Procesando reglas de la gramática desde el archivo...")

	// Procesar las reglas leídas desde el archivo
	rules := make(map[string][][]string)
	for _, line := range data {
		first, rest, separatorFound := strings.Cut(line, "->")
		if !separatorFound {
			fmt.Fprintf(os.Stdout, "⚠️  Regla \"%s\" tiene formato inválido. Recuerda agregar '->'\n", line)
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
			fmt.Fprintf(os.Stdout, "✅  Agregando regla: %s -> %s\n", first, trimmed)
		}
	}

	// Crear la estructura Grammar
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

	grammar := Grammar{
		Productions: rules,
		Terminals:   terminals,
		Initial:     "S", // Simbolo inicial por defecto
	}

	// Convertir la gramática a CNF
	chomsky := from_cfg_to_cnf(&grammar)
	fmt.Printf("🚀 CNF final: %v\n", chomsky)

	// Pedir al usuario que ingrese una frase
	fmt.Println("\n💬 Ingrese la frase que desea verificar:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sentence := strings.Split(scanner.Text(), " ")

	// Aplicar el algoritmo CYK
	accepted, table := cykParse(chomsky.Productions, sentence)
	if accepted {
		fmt.Println("✅ La frase es aceptada.")
		tree := generateParseTree(table, chomsky.Productions, sentence, chomsky.Initial)
		printTree(tree, 0)

		// Guardar el árbol como un archivo JSON
		jsonPath := "output/tree.json"
		if err := saveTreeAsJSON(tree, jsonPath); err != nil {
			fmt.Printf("⚠️  Error al guardar el árbol en JSON: %v\n", err)
		} else {
			fmt.Printf("🌳 Árbol guardado correctamente en: %s\n", jsonPath)
		}
	} else {
		fmt.Println("❌ La frase no es aceptada.")
	}
	fmt.Println("════════════════════════════════════════")
}
