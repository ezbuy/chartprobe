// Copyright © 2021 ezbuy & LITB TEAM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
