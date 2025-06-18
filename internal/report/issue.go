package report

// Severity representa o nível de severidade de um problema
type Severity string

const (
	// SeverityCritical indica um problema de segurança crítico que deve ser corrigido imediatamente
	SeverityCritical Severity = "CRITICAL"

	// SeverityWarning indica um problema de segurança que deve ser revisado
	SeverityWarning Severity = "WARNING"

	// SeverityInfo indica uma informação ou sugestão de melhoria
	SeverityInfo Severity = "INFO"
)

// Issue representa um problema de segurança encontrado durante a análise
type Issue struct {
	// Category é a categoria do problema (ssh, sysctl, services, etc)
	Category string

	// Severity é o nível de severidade do problema
	Severity Severity

	// Description é a descrição do problema
	Description string

	// CurrentValue é o valor atual da configuração
	CurrentValue string

	// RecommendedValue é o valor recomendado para a configuração
	RecommendedValue string

	// FixCommand é o comando ou configuração necessária para corrigir o problema
	FixCommand string
}
