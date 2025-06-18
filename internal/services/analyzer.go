package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mairinkdev/Hardshell/internal/report"
)

// Analyzer é o analisador de serviços
type Analyzer struct {
	mountPoint string
	rules      []ServiceRule
}

// ServiceRule representa uma regra para verificação de serviço
type ServiceRule struct {
	Name        string
	Description string
	Severity    report.Severity
	CheckFunc   func(string) bool
}

// NewAnalyzer cria um novo analisador de serviços
func NewAnalyzer(mountPoint string) *Analyzer {
	return &Analyzer{
		mountPoint: mountPoint,
		rules:      getDefaultRules(),
	}
}

// Analyze analisa os serviços ativos no sistema
func (a *Analyzer) Analyze() ([]report.Issue, error) {
	// Se estiver analisando um mountPoint, não podemos verificar serviços ativos diretamente
	if a.mountPoint != "" {
		// Verificamos os serviços habilitados olhando para os symlinks em /etc/systemd/system/multi-user.target.wants/
		systemdDir := filepath.Join(a.mountPoint, "etc/systemd/system/multi-user.target.wants")
		return a.analyzeSystemdDir(systemdDir)
	}

	// Verifica serviços ativos usando systemctl (se disponível)
	if hasCommand("systemctl") {
		return a.analyzeSystemctl()
	}

	// Alternativa para sistemas sem systemd
	if hasCommand("service") {
		return a.analyzeServiceCommand()
	}

	// Se nenhum método estiver disponível, retorna uma mensagem de erro
	return nil, fmt.Errorf("não foi possível encontrar um método para verificar serviços ativos")
}

// analyzeSystemdDir analisa os serviços habilitados em um diretório systemd
func (a *Analyzer) analyzeSystemdDir(dir string) ([]report.Issue, error) {
	var issues []report.Issue

	// Verifica se o diretório existe
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("diretório systemd não encontrado: %s", dir)
	}

	// Lista os arquivos no diretório
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler diretório systemd: %w", err)
	}

	// Analisa cada serviço
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		serviceName := file.Name()
		if strings.HasSuffix(serviceName, ".service") {
			serviceName = strings.TrimSuffix(serviceName, ".service")
		}

		// Verifica cada regra
		for _, rule := range a.rules {
			if rule.CheckFunc(serviceName) {
				issues = append(issues, report.Issue{
					Category:    "services",
					Severity:    rule.Severity,
					Description: fmt.Sprintf("%s (%s está habilitado)", rule.Description, serviceName),
					FixCommand:  fmt.Sprintf("systemctl disable %s && systemctl stop %s", serviceName, serviceName),
				})
			}
		}
	}

	return issues, nil
}

// analyzeSystemctl analisa os serviços ativos usando systemctl
func (a *Analyzer) analyzeSystemctl() ([]report.Issue, error) {
	var issues []report.Issue

	// Executa systemctl para listar serviços ativos
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--state=active", "--no-pager", "--plain", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar systemctl: %w", err)
	}

	// Processa a saída
	services := strings.Split(string(output), "\n")
	for _, service := range services {
		if service == "" {
			continue
		}

		// Extrai o nome do serviço (primeiro campo)
		fields := strings.Fields(service)
		if len(fields) == 0 {
			continue
		}

		serviceName := fields[0]
		if strings.HasSuffix(serviceName, ".service") {
			serviceName = strings.TrimSuffix(serviceName, ".service")
		}

		// Verifica cada regra
		for _, rule := range a.rules {
			if rule.CheckFunc(serviceName) {
				issues = append(issues, report.Issue{
					Category:    "services",
					Severity:    rule.Severity,
					Description: fmt.Sprintf("%s (%s está ativo)", rule.Description, serviceName),
					FixCommand:  fmt.Sprintf("systemctl disable %s && systemctl stop %s", serviceName, serviceName),
				})
			}
		}
	}

	return issues, nil
}

// analyzeServiceCommand analisa os serviços ativos usando o comando service (para sistemas sem systemd)
func (a *Analyzer) analyzeServiceCommand() ([]report.Issue, error) {
	var issues []report.Issue

	// Verifica os diretórios de init scripts
	initDirs := []string{"/etc/init.d", "/etc/rc.d"}

	for _, dir := range initDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Lista os arquivos no diretório
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler diretório de serviços: %w", err)
		}

		// Verifica cada serviço
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			serviceName := file.Name()

			// Verifica se o serviço está ativo
			cmd := exec.Command("service", serviceName, "status")
			output, _ := cmd.CombinedOutput()

			if !strings.Contains(strings.ToLower(string(output)), "stopped") &&
			   !strings.Contains(strings.ToLower(string(output)), "not running") {
				// Verifica cada regra
				for _, rule := range a.rules {
					if rule.CheckFunc(serviceName) {
						issues = append(issues, report.Issue{
							Category:    "services",
							Severity:    rule.Severity,
							Description: fmt.Sprintf("%s (%s está ativo)", rule.Description, serviceName),
							FixCommand:  fmt.Sprintf("service %s stop && update-rc.d %s disable", serviceName, serviceName),
						})
					}
				}
			}
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
		fmt.Println("Nenhum serviço inseguro encontrado.")
		return nil
	}

	// Aplica as correções
	for _, issue := range issues {
		cmd := issue.FixCommand

		// Executa o comando
		fmt.Printf("Aplicando correção: %s\n", cmd)

		// Aqui você pode implementar a execução do comando
		// Por enquanto, apenas simula a execução
		fmt.Printf("  [Simulando] %s\n", cmd)
	}

	return nil
}

// hasCommand verifica se um comando está disponível no sistema
func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// getDefaultRules retorna as regras padrão para verificação de serviços
func getDefaultRules() []ServiceRule {
	return []ServiceRule{
		{
			Name:        "telnet",
			Description: "Serviço Telnet oferece comunicação não criptografada",
			Severity:    report.SeverityCritical,
			CheckFunc: func(service string) bool {
				return service == "telnet" || service == "telnetd" || service == "inetd"
			},
		},
		{
			Name:        "rsh",
			Description: "Serviço RSH é inseguro e deve ser desabilitado",
			Severity:    report.SeverityCritical,
			CheckFunc: func(service string) bool {
				return service == "rsh" || service == "rsh-server" || service == "rshd"
			},
		},
		{
			Name:        "rlogin",
			Description: "Serviço RLogin é inseguro e deve ser desabilitado",
			Severity:    report.SeverityCritical,
			CheckFunc: func(service string) bool {
				return service == "rlogin" || service == "rlogind"
			},
		},
		{
			Name:        "rexec",
			Description: "Serviço RExec é inseguro e deve ser desabilitado",
			Severity:    report.SeverityCritical,
			CheckFunc: func(service string) bool {
				return service == "rexec" || service == "rexecd"
			},
		},
		{
			Name:        "ftp",
			Description: "Serviço FTP transfere credenciais sem criptografia",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "ftp" || service == "ftpd" ||
				       service == "vsftpd" || service == "proftpd" ||
				       strings.Contains(service, "ftp")
			},
		},
		{
			Name:        "tftp",
			Description: "Serviço TFTP é inseguro e não deve ser usado em produção",
			Severity:    report.SeverityCritical,
			CheckFunc: func(service string) bool {
				return service == "tftp" || service == "tftpd" ||
				       service == "atftpd" || service == "tftpd-hpa"
			},
		},
		{
			Name:        "finger",
			Description: "Serviço Finger pode revelar informações sobre usuários",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "finger" || service == "fingerd"
			},
		},
		{
			Name:        "talk",
			Description: "Serviço Talk não é criptografado e é raramente usado",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "talk" || service == "talkd" ||
				       service == "ntalk" || service == "ntalkd"
			},
		},
		{
			Name:        "nis",
			Description: "NIS é considerado inseguro para autenticação",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "nis" || service == "yp" ||
				       strings.Contains(service, "ypserv") ||
				       strings.Contains(service, "ypbind")
			},
		},
		{
			Name:        "snmpd",
			Description: "SNMP pode estar configurado com community strings padrão",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "snmpd" || service == "snmp"
			},
		},
		{
			Name:        "portmap",
			Description: "Portmap/RPC pode expor serviços desnecessários",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "portmap" || service == "rpcbind"
			},
		},
		{
			Name:        "sendmail",
			Description: "Servidores de email podem precisar de configuração adicional de segurança",
			Severity:    report.SeverityInfo,
			CheckFunc: func(service string) bool {
				return service == "sendmail" || service == "postfix" ||
				       service == "exim" || service == "exim4"
			},
		},
		{
			Name:        "xserver",
			Description: "Servidor X não deve estar executando em servidores",
			Severity:    report.SeverityWarning,
			CheckFunc: func(service string) bool {
				return service == "xorg" || service == "xserver" ||
				       service == "xorg-server" || service == "gdm" ||
				       service == "lightdm" || service == "kdm" ||
				       service == "sddm"
			},
		},
		{
			Name:        "avahi",
			Description: "Avahi (mDNS) não é necessário em servidores",
			Severity:    report.SeverityInfo,
			CheckFunc: func(service string) bool {
				return service == "avahi" || service == "avahi-daemon"
			},
		},
		{
			Name:        "cups",
			Description: "CUPS não é necessário em servidores sem impressoras",
			Severity:    report.SeverityInfo,
			CheckFunc: func(service string) bool {
				return service == "cups" || service == "cupsd"
			},
		},
		{
			Name:        "dhcpd",
			Description: "Servidor DHCP pode precisar de revisão de segurança",
			Severity:    report.SeverityInfo,
			CheckFunc: func(service string) bool {
				return service == "dhcpd" || service == "isc-dhcp-server"
			},
		},
	}
}
