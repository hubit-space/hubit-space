package config

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func OpenServiceConnection() (*resty.Client, error) {
	client := resty.New().
		SetTimeout(300*time.Second).
		SetBaseURL(GetEnv("BASE_URL_SERVICE", "http://localhost:8084")).
		SetHeader("Security-Code", GetEnv("SECURITY_CODE_HUBIT", "91b637d8fcd2c6da6359e6963113a1170de795e4b725b84d1e0b4ck19skdj8fm"))

	return client, nil
}
