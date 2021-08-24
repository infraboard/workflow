package gitlab

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewSCM(addr, token string) *SCM {
	return &SCM{
		Address:      addr,
		PrivateToken: token,
		Version:      "v4",
		client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type SCM struct {
	Address      string
	PrivateToken string
	Version      string

	client *http.Client
}

func (r *SCM) newJSONRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r SCM) newFormReqeust(method, url string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r *SCM) resourceURL(resource string, params map[string]string) string {
	val := make(url.Values)

	for k, v := range params {
		val.Set(k, v)
	}

	return fmt.Sprintf("%s/api/%s/%s?%s", r.Address, r.Version, resource, val.Encode())
}
