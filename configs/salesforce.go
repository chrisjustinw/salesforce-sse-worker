package configs

import "github.com/kelseyhightower/envconfig"

type SalesforceConfig struct {
	Host  string `envconfig:"HOST"`
	OrgId string `envconfig:"ORG_ID"`
}

func NewSalesforceConfig(e EnvFileRead) (SalesforceConfig, error) {
	var cfg SalesforceConfig
	if err := envconfig.Process("SALESFORCE", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
