package fixer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mairinkdev/Hardshell/internal/report"
)

// Generator é responsável pela geração de scripts de correção
type Generator struct {
	mountPoint string
}

// NewGenerator cria um novo gerador de scripts de correção
func NewGenerator(mountPoint string) *Generator {
	return &Generator{
		mountPoint: mountPoint,
	}
}

// GenerateScript gera um script shell para corrigir os problemas encontrados
func (g *Generator) GenerateScript(issues []report.Issue, outputPath string) error {
	if len(issues) == 0 {
		return fmt.Errorf("nenhum problema encontrado para corrigir")
	}

	// Prepara o conteúdo do script
	var sb strings.Builder

	// Adiciona o cabeçalho do script
	sb.WriteString("#!/bin/bash\n\n")
	sb.WriteString("# Script de correção gerado pelo Hardshell\n")
	sb.WriteString(fmt.Sprintf("# Data: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Adiciona funções de suporte
	sb.WriteString(`
# Função para exibir mensagens coloridas
function log() {
    local level=$1
    local msg=$2

    case $level in
        "INFO")
            echo -e "\033[0;34m[INFO]\033[0m $msg"
            ;;
        "WARNING")
            echo -e "\033[0;33m[WARNING]\033[0m $msg"
            ;;
        "ERROR")
            echo -e "\033[0;31m[ERROR]\033[0m $msg"
            ;;
        "SUCCESS")
            echo -e "\033[0;32m[SUCCESS]\033[0m $msg"
            ;;
        *)
            echo -e "$msg"
            ;;
    esac
}

# Função para criar backup de um arquivo
function backup_file() {
    local file=$1
    local backup="${file}.bak.$(date +%Y%m%d%H%M%S)"

    if [ -f "$file" ]; then
        cp "$file" "$backup"
        if [ $? -eq 0 ]; then
            log "INFO" "Backup criado: $backup"
            return 0
        else
            log "ERROR" "Falha ao criar backup de $file"
            return 1
        fi
    else
        log "WARNING" "Arquivo não encontrado: $file"
        return 1
    fi
}

# Verifica se o script está sendo executado como root
if [ "$EUID" -ne 0 ]; then
    log "ERROR" "Este script precisa ser executado como root."
    exit 1
fi

log "INFO" "Iniciando aplicação de correções de segurança..."

`)

	// Agrupa as issues por categoria
	categories := make(map[string][]report.Issue)
	for _, issue := range issues {
		categories[issue.Category] = append(categories[issue.Category], issue)
	}

	// Adiciona as correções para cada categoria
	for category, categoryIssues := range categories {
		sb.WriteString(fmt.Sprintf("\n# Correções para %s\n", strings.ToUpper(category)))
		sb.WriteString(fmt.Sprintf("log \"INFO\" \"Aplicando correções para %s...\"\n\n", strings.ToUpper(category)))

		// Adiciona comandos de backup para arquivos específicos
		switch category {
		case "ssh":
			sb.WriteString("# Backup do arquivo de configuração SSH\n")
			sb.WriteString("backup_file /etc/ssh/sshd_config\n\n")
		case "sysctl":
			sb.WriteString("# Backup do arquivo de configuração sysctl\n")
			sb.WriteString("backup_file /etc/sysctl.conf\n\n")
		}

		// Adiciona os comandos de correção
		for _, issue := range categoryIssues {
			sb.WriteString(fmt.Sprintf("# %s (%s)\n", issue.Description, issue.Severity))
			sb.WriteString(fmt.Sprintf("log \"INFO\" \"Corrigindo: %s\"\n", issue.Description))

			// Adapta o comando para o mountPoint, se necessário
			command := issue.FixCommand
			if g.mountPoint != "" {
				command = g.adaptCommandForMountPoint(command)
			}

			sb.WriteString(fmt.Sprintf("%s\n", command))
			sb.WriteString("if [ $? -eq 0 ]; then\n")
			sb.WriteString(fmt.Sprintf("    log \"SUCCESS\" \"Correção aplicada com sucesso: %s\"\n", issue.Description))
			sb.WriteString("else\n")
			sb.WriteString(fmt.Sprintf("    log \"ERROR\" \"Falha ao aplicar correção: %s\"\n", issue.Description))
			sb.WriteString("fi\n\n")
		}
	}

	// Adiciona o rodapé do script
	sb.WriteString("\nlog \"INFO\" \"Todas as correções foram aplicadas.\"\n")
	sb.WriteString("log \"WARNING\" \"Reinicie os serviços ou o sistema para aplicar todas as alterações.\"\n")

	// Escreve o script no arquivo
	err := os.WriteFile(outputPath, []byte(sb.String()), 0755)
	if err != nil {
		return fmt.Errorf("erro ao escrever o script de correção: %w", err)
	}

	return nil
}

// adaptCommandForMountPoint adapta um comando para uso com um ponto de montagem
func (g *Generator) adaptCommandForMountPoint(command string) string {
	// Substitui referências a arquivos específicos
	command = strings.Replace(command, "/etc/ssh/sshd_config", filepath.Join(g.mountPoint, "etc/ssh/sshd_config"), -1)
	command = strings.Replace(command, "/etc/sysctl.conf", filepath.Join(g.mountPoint, "etc/sysctl.conf"), -1)

	// Remove comandos que não podem ser executados em um mountPoint
	command = strings.Replace(command, "systemctl ", "# systemctl ", -1)
	command = strings.Replace(command, "service ", "# service ", -1)
	command = strings.Replace(command, "sysctl -p", "# sysctl -p", -1)

	return command
}
