package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	applyFixes  bool
	mountPoint  string
	outputFormat string
)

// rootCmd representa o comando base quando chamado sem subcomandos
var rootCmd = &cobra.Command{
	Use:   "hardshell",
	Short: "Ferramenta de hardening automatizado para sistemas Linux",
	Long: `Hardshell é uma ferramenta CLI para hardening automatizado de sistemas Linux.
Ela é voltada para DevOps, SREs e profissionais de segurança que precisam
garantir que os sistemas sigam as melhores práticas de segurança.

Hardshell realiza:
  - Análise de configurações de SSH
  - Verificação de configurações sysctl sensíveis à segurança
  - Identificação de serviços perigosos ativos
  - Geração de relatórios detalhados
  - Sugestão de correções através de scripts`,
}

// Execute adiciona todos os comandos filhos ao comando root e configura flags apropriadamente.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (padrão: $HOME/.hardshell.yaml)")
	rootCmd.PersistentFlags().BoolVar(&applyFixes, "apply", false, "aplicar correções automaticamente (com backup)")
	rootCmd.PersistentFlags().StringVar(&mountPoint, "mount", "", "ponto de montagem para análise (ex: /mnt)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "formato de saída (text, json, html)")
}
