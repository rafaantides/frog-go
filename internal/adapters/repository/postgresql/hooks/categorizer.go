package hooks

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

type Categorizer struct {
	categoryMap map[string]string
}

func NewCategorizer(path string) (*Categorizer, error) {
	categoryMap, err := loadCategoriesFromFile(path + "categories.json")
	if err != nil {
		return nil, err
	}
	return &Categorizer{
		categoryMap: categoryMap,
	}, nil
}

func (c *Categorizer) Categorize(name string) string {
	// Remove sucixo do tipo " - Parcela 2/3"
	parcelaRegex := regexp.MustCompile(`(?i)\s*-\s*Parcela\s+\d+/\d+$`)
	cleanName := strings.TrimSpace(parcelaRegex.ReplaceAllString(name, ""))

	if category, exists := c.categoryMap[cleanName]; exists {
		return category
	}
	return "Sem categoria"
}

func loadCategoriesFromFile(path string) (map[string]string, error) {
	if path == "" {
		return nil, nil
	}

	raw := map[string][]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	for category, names := range raw {
		for _, name := range names {
			result[name] = category
		}
	}
	return result, nil
}
