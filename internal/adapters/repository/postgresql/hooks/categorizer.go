package hooks

import (
	"encoding/json"
	"os"
	"strings"
)

type Categorizer struct {
	categories map[string][]string
}

func NewCategorizer(path string) (*Categorizer, error) {
	categoryMap, err := loadCategoriesFromFile(path + "categories.json")
	if err != nil {
		return nil, err
	}
	return &Categorizer{
		categories: categoryMap,
	}, nil
}

func (c *Categorizer) Categorize(name string) string {
	nameLower := strings.ToLower(strings.TrimSpace(name))

	for category, keywords := range c.categories {
		for _, keyword := range keywords {
			if strings.Contains(nameLower, strings.ToLower(keyword)) {
				return category
			}
		}
	}

	return "Sem categoria"
}

func loadCategoriesFromFile(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var categories map[string][]string
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
