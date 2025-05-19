package config

type Config struct {
	DNS              string   `mapstructure:"DNS" json:"DNS" yaml:"DNS"`
	Database         database `mapstructure:"DATABASE" json:"DATABASE" yaml:"DATABASE"`
	Port             string   `mapstructure:"PORT" yaml:"PORT" json:"PORT" default:"3000"`
	GoogleMapsApiKey string   `mapstructure:"GOOGLE_MAPS_API_KEY" json:"GOOGLE_MAPS_API_KEY" yaml:"GOOGLE_MAPS_API_KEY"`
	Mailer           mailer   `mapstructure:"MAILER" json:"MAILER" yaml:"MAILER"`
}

type database struct {
	Name     string `mapstructure:"NAME" yaml:"NAME"`
	Host     string `mapstructure:"HOST" yaml:"HOST"`
	User     string `mapstructure:"USER" yaml:"USER"`
	Password string `mapstructure:"PASSWORD" yaml:"PASSWORD"`
}

type mailer struct {
	Host     string `mapstructure:"HOST" json:"HOST" yaml:"HOST"`
	Port     int    `mapstructure:"PORT" json:"PORT" yaml:"PORT"`
	Username string `mapstructure:"USERNAME" json:"USERNAME" yaml:"USERNAME"`
	From     string `mapstructure:"FROM" yaml:"FROM"`
	SMTP     string `mapstructure:"SMTP" yaml:"SMTP"`
	Password string `mapstructure:"PASSWORD" yaml:"PASSWORD"`
}
