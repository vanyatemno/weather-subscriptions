package config

type Config struct {
	DNS              string   `mapstructure:"DNS" json:"DNS" yaml:"DNS"`
	Database         Database `mapstructure:"DATABASE" json:"DATABASE" yaml:"DATABASE"`
	Port             string   `mapstructure:"PORT" yaml:"PORT" json:"PORT" default:"3000"`
	FrontendURL      string   `mapstructure:"FRONTEND_URL" yaml:"FRONTEND_URL" default:"localhost:8080"`
	GoogleMapsAPIKey string   `mapstructure:"GOOGLE_MAPS_API_KEY" json:"GOOGLE_MAPS_API_KEY" yaml:"GOOGLE_MAPS_API_KEY"`
	Mailer           Mailer   `mapstructure:"MAILER" json:"MAILER" yaml:"MAILER"`
	OpenAI           OpenAI   `mapstructure:"OPENAI" json:"OPENAI" yaml:"OPENAI"`
}

type Database struct {
	Name     string `mapstructure:"NAME" yaml:"NAME"`
	Host     string `mapstructure:"HOST" yaml:"HOST"`
	User     string `mapstructure:"USER" yaml:"USER"`
	Password string `mapstructure:"PASSWORD" yaml:"PASSWORD"`
}

type Mailer struct {
	Host     string `mapstructure:"HOST" json:"HOST" yaml:"HOST"`
	Port     int    `mapstructure:"PORT" json:"PORT" yaml:"PORT"`
	Username string `mapstructure:"USERNAME" json:"USERNAME" yaml:"USERNAME"`
	From     string `mapstructure:"FROM" yaml:"FROM"`
	SMTP     string `mapstructure:"SMTP" yaml:"SMTP"`
	Password string `mapstructure:"PASSWORD" yaml:"PASSWORD"`
}

type OpenAI struct {
	OpenrouterAPIKey string `mapstructure:"OPENROUTER_API_KEY" json:"OPENROUTER_API_KEY" yaml:"OPENROUTER_API_KEY-"`
}
