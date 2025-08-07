package config

import "yatori-go-coreV3/models/ctype"

type AiSetting struct {
	AiType ctype.AiType `json:"aiType"`
	AiUrl  string       `json:"aiUrl"`
	Model  string       `json:"model"`
	APIKEY string       `json:"API_KEY" yaml:"API_KEY" mapstructure:"API_KEY"`
}
