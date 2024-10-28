package main

import (
	"fmt"
	"time"
)

// cykParse ejecuta el algoritmo CYK y mide el tiempo de ejecución.
func cykParse(grammar map[string][][]string, sentence []string) (bool, [][][]string) {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("🔍  Algoritmo CYK - Análisis de la frase")
	fmt.Println("════════════════════════════════════════")

	start := time.Now() // Inicio de la medición de tiempo

	n := len(sentence)
	T := make([][][]string, n)
	for i := range T {
		T[i] = make([][]string, n)
	}

	// Inicialización
	for j := 0; j < n; j++ {
		for lhs, productions := range grammar {
			for _, rhs := range productions {
				if len(rhs) == 1 && rhs[0] == sentence[j] {
					T[j][j] = append(T[j][j], lhs)
					fmt.Printf("📝 Terminal encontrado: %s -> %s en T[%d][%d]\n", lhs, rhs[0], j, j)
				}
			}
		}
	}

	// Llenado de la tabla
	for span := 2; span <= n; span++ {
		for i := 0; i <= n-span; i++ {
			j := i + span - 1
			for k := i; k < j; k++ {
				for lhs, productions := range grammar {
					for _, rhs := range productions {
						if len(rhs) == 2 {
							B, C := rhs[0], rhs[1]
							if contains(T[i][k], B) && contains(T[k+1][j], C) {
								fmt.Printf("🔗 Combinando: %s -> %s %s en T[%d][%d] (de T[%d][%d] y T[%d][%d])\n", lhs, B, C, i, j, i, k, k+1, j)
								T[i][j] = append(T[i][j], lhs)
							}
						}
					}
				}
			}
		}
	}

	// Depuración: Imprimir la tabla T
	printTable(T, sentence)

	// Verificar si la frase es aceptada
	accepted := contains(T[0][n-1], "S'") // Verifica con el símbolo inicial modificado

	// Fin de la medición de tiempo
	elapsed := time.Since(start)
	fmt.Printf("⏱️  Tiempo de ejecución del algoritmo CYK: %s\n", elapsed)
	fmt.Println("════════════════════════════════════════")

	return accepted, T
}

// printTable imprime la tabla de CYK para fines de depuración.
func printTable(T [][][]string, sentence []string) {
	fmt.Println("\n╔══════════════════════════════════╗")
	fmt.Println("║           Tabla de CYK           ║")
	fmt.Println("╚══════════════════════════════════╝")
	n := len(sentence)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			fmt.Printf("T[%d][%d]: %v\n", i, j, T[i][j])
		}
	}
	fmt.Println("════════════════════════════════════════")
}

// contains verifica si un slice contiene un elemento específico.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
