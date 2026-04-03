package policy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/kinghanzala/gcpsec/internal/model"
)

type Config struct {
	Rules map[string]Rule `yaml:"rules"`
}

type Rule struct {
	Severity   model.Severity `yaml:"severity"`
	Exceptions []string       `yaml:"exceptions"`
}

func DefaultConfig() Config {
	return Config{
		Rules: map[string]Rule{
			"public_buckets":     {Severity: model.SeverityHigh},
			"owner_roles":        {Severity: model.SeverityHigh},
			"public_iam_members": {Severity: model.SeverityHigh},
			"open_ssh":           {Severity: model.SeverityMedium},
			"open_rdp":           {Severity: model.SeverityMedium},
		},
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, fmt.Errorf("failed to read config %q: %w", cleanPath, err)
	}

	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config %q: %w", cleanPath, err)
	}

	for key, rule := range fileCfg.Rules {
		defaultRule, ok := cfg.Rules[key]
		if !ok {
			cfg.Rules[key] = rule
			continue
		}

		if rule.Severity != "" {
			defaultRule.Severity = rule.Severity
		}
		if len(rule.Exceptions) > 0 {
			defaultRule.Exceptions = append([]string(nil), rule.Exceptions...)
		}
		cfg.Rules[key] = defaultRule
	}

	return cfg, nil
}
