package gitlab

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewRepository(addr, token string) *Repository {
	return &Repository{
		Address:      addr,
		PrivateToken: token,
		Version:      "v4",
		client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type Repository struct {
	Address      string
	PrivateToken string
	Version      string

	client *http.Client
}

func (r *Repository) newJSONRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r Repository) newFormReqeust(method, url string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r *Repository) resourceURL(resource string, params map[string]string) string {
	val := make(url.Values)

	for k, v := range params {
		val.Set(k, v)
	}

	return fmt.Sprintf("%s/api/%s/%s?%s", r.Address, r.Version, resource, val.Encode())
}
