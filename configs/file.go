package configs

import (
	"github.com/joho/godotenv"
)

type EnvFileRead struct{}

func ReadEnvFile() (EnvFileRead, error) {
	return EnvFileRead{}, godotenv.Load(".env")
}
