package client

import (
	"errors"
	"strings"
)

type Chart struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var ErrInvalidChartFormation = errors.New("chart: invalid formation")

func NewChart(cname string, version string) *Chart {
	return &Chart{
		Name:    cname,
		Version: version,
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
