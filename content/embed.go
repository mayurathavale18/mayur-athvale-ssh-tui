package content

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed portfolio.yaml
var portfolioData []byte

func LoadPortfolio() (Portfolio, error) {
	var p Portfolio
	if err := yaml.Unmarshal(portfolioData, &p); err != nil {
		return Portfolio{}, err
	}
	return p, nil
}
