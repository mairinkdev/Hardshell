package ssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mairinkdev/Hardshell/internal/report"
)

// Analyzer é o analisador de configurações SSH
type Analyzer struct {
	mountPoint string
	configPath string
	rules      []SSHRule
}

// SSHRule representa uma regra para verificação de configuração SSH
type SSHRule struct {
	Key               string
	RecommendedValue  string
	Severity          report.Severity
	Description       string
	ComparisonFunc    func(string, string) bool
}

// NewAnalyzer cria um novo analisador SSH
func NewAnalyzer(mountPoint string) *Analyzer {
	configPath := "/etc/ssh/sshd_config"
	if mountPoint != "" {
		configPath = filepath.Join(mountPoint, configPath)
	}

	return &Analyzer{
		mountPoint: mountPoint,
		configPath: configPath,
		rules:      getDefaultRules(),
	}
}

// Analyze analisa o arquivo sshd_config em busca de configurações inseguras
func (a *Analyzer) Analyze() ([]report.Issue, error) {
	// Verifica se o arquivo de configuração existe
	if _, err := os.Stat(a.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração SSH não encontrado: %s", a.configPath)
	}

	// Lê o arquivo de configuração
	configFile, err := os.Open(a.configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de configuração SSH: %w", err)
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
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração SSH: %w", err)
	}

	// Verifica as regras
	var issues []report.Issue

	for _, rule := range a.rules {
		value, exists := config[rule.Key]

		// Se a configuração não existe, considere como uma violação
		if !exists {
			issues = append(issues, report.Issue{
				Category:         "ssh",
				Severity:         rule.Severity,
				Description:      rule.Description,
				RecommendedValue: rule.RecommendedValue,
				FixCommand:       fmt.Sprintf("echo '%s %s' >> /etc/ssh/sshd_config", rule.Key, rule.RecommendedValue),
			})
			continue
		}

		// Verifica se o valor atual atende à regra
		if !rule.ComparisonFunc(value, rule.RecommendedValue) {
			issues = append(issues, report.Issue{
				Category:         "ssh",
				Severity:         rule.Severity,
				Description:      rule.Description,
				CurrentValue:     value,
				RecommendedValue: rule.RecommendedValue,
				FixCommand:       fmt.Sprintf("sed -i 's/^%s.*/%s %s/' /etc/ssh/sshd_config", rule.Key, rule.Key, rule.RecommendedValue),
			})
		}
	}

	return issues, nil
}

// Fix gera um script para corrigir os problemas encontrados
func (a *Analyzer) Fix() error {
	// Analisa os problemas
	issues, err := a.Analyze()
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		fmt.Println("Nenhum problema encontrado nas configurações SSH.")
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
			cmd = strings.Replace(cmd, "/etc/ssh/sshd_config", a.configPath, -1)
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

// getDefaultRules retorna as regras padrão para verificação SSH
func getDefaultRules() []SSHRule {
	return []SSHRule{
		{
			Key:              "PermitRootLogin",
			RecommendedValue: "no",
			Severity:         report.SeverityCritical,
			Description:      "Login direto como root deve ser desabilitado",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "Protocol",
			RecommendedValue: "2",
			Severity:         report.SeverityCritical,
			Description:      "Apenas o protocolo SSH 2 deve ser permitido (SSH 1 é inseguro)",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "PasswordAuthentication",
			RecommendedValue: "no",
			Severity:         report.SeverityWarning,
			Description:      "Autenticação por senha deve ser desabilitada, prefira chaves SSH",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "PermitEmptyPasswords",
			RecommendedValue: "no",
			Severity:         report.SeverityCritical,
			Description:      "Senhas vazias não devem ser permitidas",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "X11Forwarding",
			RecommendedValue: "no",
			Severity:         report.SeverityWarning,
			Description:      "Encaminhamento X11 deve ser desabilitado se não for necessário",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
		{
			Key:              "MaxAuthTries",
			RecommendedValue: "4",
			Severity:         report.SeverityWarning,
			Description:      "Número máximo de tentativas de autenticação deve ser limitado",
			ComparisonFunc:   func(actual, recommended string) bool {
				// Verifica se o valor atual é menor ou igual ao recomendado
				return actual <= recommended
			},
		},
		{
			Key:              "ClientAliveInterval",
			RecommendedValue: "300",
			Severity:         report.SeverityInfo,
			Description:      "Definir um intervalo de keepalive para detectar clientes desconectados",
			ComparisonFunc:   func(actual, recommended string) bool {
				// Verifica se há algum valor definido
				return actual != ""
			},
		},
		{
			Key:              "ClientAliveCountMax",
			RecommendedValue: "3",
			Severity:         report.SeverityInfo,
			Description:      "Limitar o número de mensagens keepalive sem resposta antes de desconectar",
			ComparisonFunc:   func(actual, recommended string) bool {
				// Verifica se o valor atual é menor ou igual ao recomendado
				return actual <= recommended
			},
		},
		{
			Key:              "LogLevel",
			RecommendedValue: "VERBOSE",
			Severity:         report.SeverityWarning,
			Description:      "Nível de log deve ser detalhado para auditoria adequada",
			ComparisonFunc:   func(actual, recommended string) bool {
				// Valores aceitáveis: VERBOSE ou INFO
				return actual == "VERBOSE" || actual == "INFO"
			},
		},
		{
			Key:              "UsePAM",
			RecommendedValue: "yes",
			Severity:         report.SeverityWarning,
			Description:      "PAM deve ser habilitado para controle de acesso avançado",
			ComparisonFunc:   func(actual, recommended string) bool { return actual == recommended },
		},
	}
}
