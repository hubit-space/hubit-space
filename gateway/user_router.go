package repository

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

type DocNumRepositoryClient struct {
	client *resty.Client
}

type DocNumRepository interface {
	GetDocNum(docname string, year int, month int, day int) (map[string]any, error)
}

func NewDocNumRepository(client *resty.Client) *DocNumRepositoryClient {
	return &DocNumRepositoryClient{
		client: client,
	}
}

func (r *DocNumRepositoryClient) GetDocNum(docname string, year int, month int, day int) (map[string]any, error) {
	var response map[string]any

	payload := map[string]any{
		"docname": docname,
		"year":    year,
		"month":   month,
		"day":     day,
	}

	resp, err := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&response).
		SetBody(payload).
		Post(os.Getenv("BASE_URL_SERVICE") + "/docnum/get-last")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get doc num")
	}

	return response, nil
}
