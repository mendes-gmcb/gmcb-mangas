package env

var DB_CONFIG dbConfig
var AWS_CONFIG awsConfig

type dbConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Pass     string `env:"DB_PASS"`
	Database string `env:"DB_NAME"`
}

type awsConfig struct {
	AcessKey  string `env:"AWS_ACCESS_KEY"`
	SecretKey string `env:"AWS_SECRET_KEY"`
	Region    string `env:"AWS_REGION"`
}
