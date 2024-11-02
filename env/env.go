package env

import (
	"os"
	"strconv"
)

func GetEnv(valor string) string {
	return os.Getenv(valor)
}

func GetEnvInt32(valor string) int32 {
	value := os.Getenv(valor)
	v, _ := strconv.ParseInt(value, 10, 32)
	return int32(v)
}

func GetEnvBool(valor string) bool {
	return os.Getenv(valor) == "true"
}