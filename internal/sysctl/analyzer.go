package sysctl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mairinkdev/Hardshell/internal/report"
)

// Analyzer é o analisador de configurações sysctl
type Analyzer struct {
	mountPoint string
	configPath string
	rules      []SysctlRule
}

// SysctlRule representa uma regra para verificação de configuração sysctl
type SysctlRule struct {
	Key               string
	RecommendedValue  string
	Severity          report.Severity
	Description       string
	ComparisonFunc    func(string, string) bool
}

// NewAnalyzer cria um novo analisador sysctl
func NewAnalyzer(mountPoint string) *Analyzer {
	configPath := "/etc/sysctl.conf"
	if mountPoint != "" {
		configPath = filepath.Join(mountPoint, configPath)
	}

	return &Analyzer{
		mountPoint: mountPoint,
		configPath: configPath,
		rules:      getDefaultRules(),
	}
}

// Analyze analisa as configurações sysctl relacionadas à segurança
func (a *Analyzer) Analyze() ([]report.Issue, error) {
	// Verifica se o arquivo de configuração existe
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração sysctl não encontrado: %s", a.configPath)
	}

	// Lê o arquivo de configuração
	configFile, err := os.Open(a.configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de configuração sysctl: %w", err)
	}
	defer configFile.Close()

	// Analisa o arquivo de configuração
	config := make(map[string]string)
	scanner := bufio.NewScanner(configFile)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignora comentários e linhas em branco
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Separa a chave e o valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração sysctl: %w", err)
	}

	// Verifica arquivos adicionais em sysctl.d se não estiver em um mountPoint
	if a.mountPoint == "" {
		err = a.readSysctlD(config)
		if err != nil {
			return nil, err
		}
	} else {
		// Usa o mountPoint para verificar o diretório sysctl.d
		sysctlDPath := filepath.Join(a.mountPoint, "/etc/sysctl.d")
		err = a.readSysctlDFromPath(sysctlDPath, config)
		if err != nil {
			return nil, err
		}
	}

	// Verifica as regras
	var issues []report.Issue

	for _, rule := range a.rules {
		value, exists := config[rule.Key]

		// Se a configuração não existe, considere como uma violação
		if !exists {
			issues = append(issues, report.Issue{
				Category:         "sysctl",
				Severity:         rule.Severity,
				Description:      rule.Description,
				RecommendedValue: rule.RecommendedValue,
				FixCommand:       fmt.Sprintf("echo '%s = %s' >> /etc/sysctl.conf && sysctl -p", rule.Key, rule.RecommendedValue),
			})
			continue
		}

		// Verifica se o valor atual atende à regra
		if !rule.ComparisonFunc(value, rule.RecommendedValue) {
			issues = append(issues, report.Issue{
				Category:         "sysctl",
				Severity:         rule.Severity,
				Description:      rule.Description,
				CurrentValue:     value,
				RecommendedValue: rule.RecommendedValue,
				FixCommand:       fmt.Sprintf("sed -i 's/^%s.*/%s = %s/' /etc/sysctl.conf && sysctl -p", rule.Key, rule.Key, rule.RecommendedValue),
			})
		}
	}

	return issues, nil
}

// readSysctlD lê os arquivos .conf em /etc/sysctl.d/
func (a *Analyzer) readSysctlD(config map[string]string) error {
	sysctlDPath := "/etc/sysctl.d"
	return a.readSysctlDFromPath(sysctlDPath, config)
}

// readSysctlDFromPath lê os arquivos .conf em um diretório específico
func (a *Analyzer) readSysctlDFromPath(dirPath string, config map[string]string) error {
	// Verifica se o diretório existe
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil // Não é um erro, apenas não existem arquivos adicionais
	}

	// Lê o diretório
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório sysctl.d: %w", err)
	}

	// Analisa cada arquivo .conf
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".conf") {
			filePath := filepath.Join(dirPath, file.Name())

			// Abre e lê o arquivo
			configFile, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("erro ao abrir arquivo de configuração sysctl: %w", err)
			}

			scanner := bufio.NewScanner(configFile)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())

				// Ignora comentários e linhas em branco
				if strings.HasPrefix(line, "#") || line == "" {
					continue
				}

				// Separa a chave e o valor
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					continue
				}

				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Apenas sobrescreve se não existir (arquivos em /etc/sysctl.conf têm precedência)
				if _, exists := config[key]; !exists {
					config[key] = value
				}
			}

			configFile.Close()

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("erro ao ler arquivo de configuração sysctl: %w", err)
			}
		}
	}

	return nil
}

// Fix gera um script para corrigir os problemas encontrados
func (a *Analyzer) Fix() error {
	// Analisa os problemas
	issues, err := a.Analyze()
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		fmt.Println("Nenhum problema encontrado nas configurações sysctl.")
		return nil
	}

	// Cria um backup do arquivo de configuração
	backupPath := a.configPath + ".bak"
	err = copyFile(a.configPath, backupPath)
	if err != nil {
		return fmt.Errorf("erro ao criar backup do arquivo de configuração: %w", err)
	}

	fmt.Printf("Backup criado em %s\n", backupPath)

	// Aplica as correções
	for _, issue := range issues {
		cmd := issue.FixCommand

		// Adapta o comando para o mountPoint, se necessário
		if a.mountPoint != "" {
			cmd = strings.Replace(cmd, "/etc/sysctl.conf", a.configPath, -1)
			// Remove o comando sysctl -p se estiver em um mountPoint
			cmd = strings.Replace(cmd, "&& sysctl -p", "", -1)
		}

		// Executa o comando
		fmt.Printf("Aplicando correção: %s\n", cmd)

		// Aqui você pode implementar a execução do comando
		// Por enquanto, apenas simula a execução
		fmt.Printf("  [Simulando] %s\n", cmd)
	}

	return nil
}

// copyFile copia um arquivo de origem para destino
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// getDefaultRules retorna as regras padrão para verificação sysctl
func getDefaultRules() []SysctlRule {
	return []SysctlRule{
		{
			Key:              "net.ipv4.tcp_syncookies",
			RecommendedValue: "1",
			Severity:         report.SeverityCritical,
			Description:      "SYN flood protection deve estar habilitada",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.conf.all.accept_redirects",
			RecommendedValue: "0",
			Severity:         report.SeverityWarning,
			Description:      "ICMP redirects não devem ser aceitos",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.conf.all.send_redirects",
			RecommendedValue: "0",
			Severity:         report.SeverityWarning,
			Description:      "ICMP redirects não devem ser enviados",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.conf.all.accept_source_route",
			RecommendedValue: "0",
			Severity:         report.SeverityCritical,
			Description:      "Source routing deve estar desabilitado",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.conf.all.log_martians",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "Pacotes com endereços impossíveis devem ser registrados",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.icmp_echo_ignore_broadcasts",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "ICMP broadcasts devem ser ignorados",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.icmp_ignore_bogus_error_responses",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "Mensagens de erro ICMP malformadas devem ser ignoradas",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "net.ipv4.tcp_rfc1337",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "Proteção contra TIME-WAIT assassination deve estar habilitada",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "kernel.randomize_va_space",
			RecommendedValue: "2",
			Severity:         report.SeverityCritical,
			Description:      "ASLR deve estar totalmente habilitado",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "fs.protected_hardlinks",
			RecommendedValue: "1",
			Severity:         report.SeverityCritical,
			Description:      "Hard links devem ser protegidos",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "fs.protected_symlinks",
			RecommendedValue: "1",
			Severity:         report.SeverityCritical,
			Description:      "Symbolic links devem ser protegidos",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "kernel.kptr_restrict",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "Restrição de exibição de ponteiros do kernel deve estar habilitada",
			ComparisonFunc:   func(actual, recommended string) bool {
				// Verifica se o valor atual é pelo menos 1
				actualInt, err := strconv.Atoi(actual)
				if err != nil {
					return false
				}
				return actualInt >= 1
			},
		},
		{
			Key:              "kernel.dmesg_restrict",
			RecommendedValue: "1",
			Severity:         report.SeverityWarning,
			Description:      "Acesso ao dmesg deve ser restrito",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "kernel.sysrq",
			RecommendedValue: "0",
			Severity:         report.SeverityWarning,
			Description:      "SysRq deve estar desabilitado em ambientes de produção",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "kernel.core_uses_pid",
			RecommendedValue: "1",
			Severity:         report.SeverityInfo,
			Description:      "Core dumps devem incluir o PID no nome do arquivo",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
	}
}
