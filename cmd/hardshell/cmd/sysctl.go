package cmd

import (
	"fmt"

	"github.com/mairinkdev/Hardshell/internal/sysctl"
	"github.com/spf13/cobra"
)

// sysctlCmd representa o comando sysctl
var sysctlCmd = &cobra.Command{
	Use:   "sysctl",
	Short: "Analisa as configurações sysctl",
	Long: `Verifica as configurações do sysctl.conf e arquivos em sysctl.d para
garantir que as configurações relacionadas à segurança estão adequadas.
Analisa parâmetros como:
  - net.ipv4.tcp_syncookies
  - net.ipv4.conf.all.accept_redirects
  - kernel.randomize_va_space
  - fs.protected_hardlinks/symlinks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Analisando configurações sysctl...")

		// Cria o analisador sysctl
		analyzer := sysctl.NewAnalyzer(mountPoint)

		// Executa a análise
		issues, err := analyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar configurações sysctl: %w", err)
		}

		// Exibe os resultados
		fmt.Printf("Encontradas %d questões nas configurações sysctl\n", len(issues))
		for _, issue := range issues {
			fmt.Printf("[%s] %s\n", issue.Severity, issue.Description)
		}

		// Se --apply foi especificado, gerar e aplicar correções
		if applyFixes {
			fmt.Println("Aplicando correções...")
			if err := analyzer.Fix(); err != nil {
				return fmt.Errorf("erro ao aplicar correções: %w", err)
			}
			fmt.Println("Correções aplicadas com sucesso!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sysctlCmd)
}
