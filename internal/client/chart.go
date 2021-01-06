package client

import (
	"errors"
	"strings"
	"time"

	"helm.sh/helm/pkg/repo"
)

type Chart struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
}

var ErrInvalidChartFormation = errors.New("chart: invalid formation")

func NewChart(c *repo.ChartVersion) *Chart {
	return &Chart{
		Name:    c.Name,
		Version: c.Version,
		Created: c.Created,
	}
}

func ParseChart(cstr string) (*Chart, error) {
	parts := strings.Split(cstr, ":")
	if len(parts) != 2 {
		return nil, ErrInvalidChartFormation
	}
	return &Chart{
		Name:    parts[0],
		Version: parts[1],
	}, nil
}
