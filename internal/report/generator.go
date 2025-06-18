package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Generator é responsável pela geração de relatórios
type Generator struct {
	format string
}

// NewGenerator cria um novo gerador de relatórios
func NewGenerator(format string) *Generator {
	return &Generator{
		format: format,
	}
}

// Generate gera um relatório baseado nas issues encontradas
func (g *Generator) Generate(issues []Issue) (string, error) {
	switch strings.ToLower(g.format) {
	case "json":
		return g.generateJSON(issues)
	case "html":
		return g.generateHTML(issues)
	default:
		return g.generateText(issues)
	}
}

// generateText gera um relatório em formato texto
func (g *Generator) generateText(issues []Issue) (string, error) {
	var sb strings.Builder

	sb.WriteString("=== RELATÓRIO DE SEGURANÇA HARDSHELL ===\n\n")

	// Agrupa as issues por categoria
	categories := make(map[string][]Issue)
	for _, issue := range issues {
		categories[issue.Category] = append(categories[issue.Category], issue)
	}

	// Gera o relatório para cada categoria
	for category, categoryIssues := range categories {
		sb.WriteString(fmt.Sprintf("== %s ==\n", strings.ToUpper(category)))

		for i, issue := range categoryIssues {
			sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, issue.Severity, issue.Description))

			if issue.CurrentValue != "" {
				sb.WriteString(fmt.Sprintf("   Valor atual: %s\n", issue.CurrentValue))
			}

			if issue.RecommendedValue != "" {
				sb.WriteString(fmt.Sprintf("   Valor recomendado: %s\n", issue.RecommendedValue))
			}

			if issue.FixCommand != "" {
				sb.WriteString(fmt.Sprintf("   Correção: %s\n", issue.FixCommand))
			}

			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// generateJSON gera um relatório em formato JSON
func (g *Generator) generateJSON(issues []Issue) (string, error) {
	type Report struct {
		Issues  []Issue `json:"issues"`
		Summary struct {
			Critical int `json:"critical"`
			Warning  int `json:"warning"`
			Info     int `json:"info"`
			Total    int `json:"total"`
		} `json:"summary"`
	}

	report := Report{
		Issues: issues,
	}

	// Calcula o resumo
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityCritical:
			report.Summary.Critical++
		case SeverityWarning:
			report.Summary.Warning++
		case SeverityInfo:
			report.Summary.Info++
		}
	}

	report.Summary.Total = len(issues)

	// Serializa para JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// generateHTML gera um relatório em formato HTML
func (g *Generator) generateHTML(issues []Issue) (string, error) {
	var sb strings.Builder

	// Cabeçalho HTML
	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Relatório de Segurança Hardshell</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            line-height: 1.6;
            max-width: 960px;
            margin: 0 auto;
            padding: 20px;
            background-color: #1e1e1e;
            color: #f0f0f0;
        }
        h1, h2 {
            color: #00cc00;
            border-bottom: 1px solid #444;
            padding-bottom: 5px;
        }
        .issue {
            margin-bottom: 20px;
            padding: 10px;
            border-left: 5px solid #555;
            background-color: #2a2a2a;
        }
        .CRITICAL {
            border-left-color: #ff3333;
        }
        .WARNING {
            border-left-color: #ffcc00;
        }
        .INFO {
            border-left-color: #3399ff;
        }
        .severity {
            font-weight: bold;
            padding: 2px 8px;
            border-radius: 3px;
            display: inline-block;
        }
        .severity.CRITICAL {
            background-color: #ff3333;
            color: white;
        }
        .severity.WARNING {
            background-color: #ffcc00;
            color: black;
        }
        .severity.INFO {
            background-color: #3399ff;
            color: white;
        }
        .summary {
            display: flex;
            margin: 20px 0;
        }
        .summary-item {
            flex: 1;
            text-align: center;
            padding: 10px;
            margin: 5px;
            background-color: #2a2a2a;
            border-radius: 5px;
        }
        .summary-item.critical {
            border-bottom: 3px solid #ff3333;
        }
        .summary-item.warning {
            border-bottom: 3px solid #ffcc00;
        }
        .summary-item.info {
            border-bottom: 3px solid #3399ff;
        }
        .summary-item.total {
            border-bottom: 3px solid #00cc00;
        }
        .summary-number {
            font-size: 2em;
            font-weight: bold;
        }
        .fix {
            background-color: #333;
            padding: 8px;
            border-radius: 3px;
            font-family: monospace;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <h1>Relatório de Segurança Hardshell</h1>
`)

	// Resumo
	var critical, warning, info int
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityCritical:
			critical++
		case SeverityWarning:
			warning++
		case SeverityInfo:
			info++
		}
	}

	sb.WriteString(`    <div class="summary">
        <div class="summary-item critical">
            <div class="summary-number">`)
	sb.WriteString(fmt.Sprintf("%d", critical))
	sb.WriteString(`</div>
            <div>Críticos</div>
        </div>
        <div class="summary-item warning">
            <div class="summary-number">`)
	sb.WriteString(fmt.Sprintf("%d", warning))
	sb.WriteString(`</div>
            <div>Avisos</div>
        </div>
        <div class="summary-item info">
            <div class="summary-number">`)
	sb.WriteString(fmt.Sprintf("%d", info))
	sb.WriteString(`</div>
            <div>Informações</div>
        </div>
        <div class="summary-item total">
            <div class="summary-number">`)
	sb.WriteString(fmt.Sprintf("%d", len(issues)))
	sb.WriteString(`</div>
            <div>Total</div>
        </div>
    </div>
`)

	// Agrupa as issues por categoria
	categories := make(map[string][]Issue)
	for _, issue := range issues {
		categories[issue.Category] = append(categories[issue.Category], issue)
	}

	// Gera o relatório para cada categoria
	for category, categoryIssues := range categories {
		sb.WriteString(fmt.Sprintf("    <h2>%s</h2>\n", strings.ToUpper(category)))

		for _, issue := range categoryIssues {
			sb.WriteString(fmt.Sprintf(`    <div class="issue %s">
        <span class="severity %s">%s</span>
        <p><strong>%s</strong></p>
`, issue.Severity, issue.Severity, issue.Severity, issue.Description))

			if issue.CurrentValue != "" {
				sb.WriteString(fmt.Sprintf("        <p>Valor atual: <code>%s</code></p>\n", issue.CurrentValue))
			}

			if issue.RecommendedValue != "" {
				sb.WriteString(fmt.Sprintf("        <p>Valor recomendado: <code>%s</code></p>\n", issue.RecommendedValue))
			}

			if issue.FixCommand != "" {
				sb.WriteString(fmt.Sprintf("        <p>Correção:</p>\n        <pre class=\"fix\">%s</pre>\n", issue.FixCommand))
			}

			sb.WriteString("    </div>\n")
		}
	}

	// Rodapé HTML
	sb.WriteString(`</body>
</html>`)

	return sb.String(), nil
}
