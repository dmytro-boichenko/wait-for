package waiter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	elasticHealthPath = "_cat/health?h=status"
)

func elasticsearchWaiter() waiter {
	host := envVar("ELASTIC_HOST", "http://localhost")
	port := envVar("ELASTIC_PORT", "9200")

	connectionString := fmt.Sprintf("%s:%s/%s", host, port, elasticHealthPath)

	return ElasticsearchWaiter{
		connectionString: connectionString,
	}
}

type ElasticsearchWaiter struct {
	connectionString string
}

func (w ElasticsearchWaiter) waitFor() (bool, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, w.connectionString, nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return false, err
	}

	s := strings.TrimSpace(string(body))

	if "green" != s && "yellow" != s {
		return false, errors.Errorf("%s health is '%s'", w.name(), s)
	}

	return true, nil
}

func (w ElasticsearchWaiter) name() string {
	return "Elasticsearch"
}
