package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mairinkdev/Hardshell/internal/sysctl"
)

func main() {
	// Obtém o diretório atual
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Erro ao obter diretório atual: %v\n", err)
		os.Exit(1)
	}

	// Cria um diretório temporário para os testes
	tempDir := filepath.Join(dir, "tests", "temp")
	os.MkdirAll(tempDir, 0755)
	os.MkdirAll(filepath.Join(tempDir, "etc"), 0755)

	// Copia o arquivo de exemplo para o local esperado
	srcFile := filepath.Join(dir, "tests", "fixtures", "sysctl.conf")
	destFile := filepath.Join(tempDir, "etc", "sysctl.conf")

	// Lê o arquivo de origem
	content, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo de exemplo: %v\n", err)
		os.Exit(1)
	}

	// Escreve o arquivo de destino
	err = os.WriteFile(destFile, content, 0644)
	if err != nil {
		fmt.Printf("Erro ao escrever arquivo temporário: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Arquivo de teste copiado para:", destFile)

	// Cria um analisador sysctl apontando para o diretório temporário
	analyzer := sysctl.NewAnalyzer(tempDir)

	// Tenta analisar o arquivo
	issues, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Erro ao analisar configuração sysctl: %v\n", err)
		os.Exit(1)
	}

	// Exibe os resultados
	fmt.Printf("Análise concluída! Encontrados %d problemas:\n", len(issues))
	for i, issue := range issues {
		fmt.Printf("%d. [%s] %s\n", i+1, issue.Severity, issue.Description)
		if issue.CurrentValue != "" {
			fmt.Printf("   Valor atual: %s\n", issue.CurrentValue)
		}
		if issue.RecommendedValue != "" {
			fmt.Printf("   Valor recomendado: %s\n", issue.RecommendedValue)
		}
		fmt.Println()
	}

	// Limpa os arquivos temporários
	os.RemoveAll(tempDir)
}
