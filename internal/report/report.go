package report

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chenpingan-dev/patch-orbit/internal/impact"
)

//go:embed web/*
var assets embed.FS

func HTML(value impact.Report) ([]byte, error) {
	template, err := assets.ReadFile("web/report.html")
	if err != nil {
		return nil, err
	}
	css, err := assets.ReadFile("web/style.css")
	if err != nil {
		return nil, err
	}
	app, err := assets.ReadFile("web/app.js")
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("report JSON: %w", err)
	}
	result := strings.Replace(string(template), "__STYLE__", string(css), 1)
	result = strings.Replace(result, "__APP__", string(app), 1)
	result = strings.Replace(result, "__REPORT__", string(payload), 1)
	return []byte(result), nil
}
