package main

import (
	"fmt"
	"os"

	"github.com/mairinkdev/Hardshell/cmd/hardshell/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar Hardshell: %s\n", err)
		os.Exit(1)
	}
}
