package cmd

import (
	"fmt"

	"github.com/mairinkdev/Hardshell/internal/report"
	"github.com/mairinkdev/Hardshell/internal/services"
	"github.com/mairinkdev/Hardshell/internal/ssh"
	"github.com/mairinkdev/Hardshell/internal/sysctl"
	"github.com/spf13/cobra"
)

// scanCmd representa o comando scan
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Realiza um scan completo do sistema",
	Long: `Executa uma verificação completa de segurança no sistema,
analisando configurações SSH, sysctl e serviços ativos.
Gera um relatório detalhado com as descobertas e recomendações.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Iniciando scan completo do sistema...")

		// Cria os analisadores
		sshAnalyzer := ssh.NewAnalyzer(mountPoint)
		sysctlAnalyzer := sysctl.NewAnalyzer(mountPoint)
		servicesAnalyzer := services.NewAnalyzer(mountPoint)

		// Executa as análises
		sshIssues, err := sshAnalyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar configurações SSH: %w", err)
		}

		sysctlIssues, err := sysctlAnalyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar configurações sysctl: %w", err)
		}

		servicesIssues, err := servicesAnalyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar serviços: %w", err)
		}

		// Cria e exibe o relatório
		reportGenerator := report.NewGenerator(outputFormat)

		allIssues := []report.Issue{}
		allIssues = append(allIssues, sshIssues...)
		allIssues = append(allIssues, sysctlIssues...)
		allIssues = append(allIssues, servicesIssues...)

		reportData, err := reportGenerator.Generate(allIssues)
		if err != nil {
			return fmt.Errorf("erro ao gerar relatório: %w", err)
		}

		fmt.Println(reportData)

		// Resumo das descobertas
		var critical, warning, info int
		for _, issue := range allIssues {
			switch issue.Severity {
			case "CRITICAL":
				critical++
			case "WARNING":
				warning++
			case "INFO":
				info++
			}
		}

		fmt.Printf("\nResumo do scan:\n")
		fmt.Printf("  Problemas críticos: %d\n", critical)
		fmt.Printf("  Avisos: %d\n", warning)
		fmt.Printf("  Informações: %d\n", info)
		fmt.Printf("  Total: %d\n", len(allIssues))

		// Se --apply foi especificado, gerar e aplicar correções
		if applyFixes {
			fmt.Println("\nAplicando correções...")

			if err := sshAnalyzer.Fix(); err != nil {
				return fmt.Errorf("erro ao aplicar correções SSH: %w", err)
			}

			if err := sysctlAnalyzer.Fix(); err != nil {
				return fmt.Errorf("erro ao aplicar correções sysctl: %w", err)
			}

			if err := servicesAnalyzer.Fix(); err != nil {
				return fmt.Errorf("erro ao aplicar correções de serviços: %w", err)
			}

			fmt.Println("Todas as correções foram aplicadas com sucesso!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
