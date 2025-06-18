package cmd

import (
	"fmt"

	"github.com/mairinkdev/Hardshell/internal/ssh"
	"github.com/spf13/cobra"
)

// sshCmd representa o comando ssh
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Analisa a configuração do SSH",
	Long: `Verifica a configuração do sshd_config em busca de configurações inseguras
como PermitRootLogin, Protocol, PasswordAuthentication e outras opções críticas.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Analisando configurações SSH...")

		// Cria o analisador SSH
		analyzer := ssh.NewAnalyzer(mountPoint)

		// Executa a análise
		issues, err := analyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar configurações SSH: %w", err)
		}

		// Exibe os resultados
		fmt.Printf("Encontradas %d questões nas configurações SSH\n", len(issues))
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
	rootCmd.AddCommand(sshCmd)
}
