package api

type RestApiConfig struct {
	BasePath string `env:"API_BASE_PATH"`
	Port     string `env:"API_PORT"`
}
