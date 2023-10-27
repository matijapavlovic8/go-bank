package postgres

type DbConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}
