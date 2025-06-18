package cmd

import (
	"fmt"

	"github.com/mairinkdev/Hardshell/internal/services"
	"github.com/spf13/cobra"
)

// servicesCmd representa o comando services
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Analisa os serviços ativos",
	Long: `Verifica os serviços ativos no sistema para identificar serviços potencialmente perigosos
ou mal configurados. Inclui verificação de serviços como telnet, rsh, rlogin, e outros
serviços inseguros.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Analisando serviços ativos...")

		// Cria o analisador de serviços
		analyzer := services.NewAnalyzer(mountPoint)

		// Executa a análise
		issues, err := analyzer.Analyze()
		if err != nil {
			return fmt.Errorf("erro ao analisar serviços: %w", err)
		}

		// Exibe os resultados
		fmt.Printf("Encontrados %d serviços com problemas de segurança\n", len(issues))
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
	rootCmd.AddCommand(servicesCmd)
}
