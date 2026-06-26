package config

type JWT struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Expire int64  `mapstructure:"expire" json:"expire" yaml:"expire"`
	Issuer string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}
