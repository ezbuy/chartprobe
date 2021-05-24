package client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"helm.sh/helm/pkg/repo"
)

var ErrPurgeConflict = errors.New("client: purge only can be executed with empty charts list")

var ClientOnce sync.Once
var defaultClient MuseumClient

func envOrYAML(field string) (string, error) {
	value := os.Getenv(field)
	if value == "" {
		var ok bool
		value, ok = viper.Get(field).(string)
		if !ok {
			return "", fmt.Errorf("conf: field[%s] not found", field)
		}
	}
	if value == "" {
		return "", fmt.Errorf("conf: field[%s] not found", field)
	}
	return value, nil
}

func NewClient() MuseumClient {
	ClientOnce.Do(func() {
		// TODO: multi-tenant
		host, err := envOrYAML("host")
		if err != nil {
			panic(err)
		}
		timeoutSecond, ok := viper.Get("timeout").(int)
		if !ok {
			timeoutSecond = 10
		}
		// TODO: auth
		c := http.Client{
			Timeout: time.Duration(timeoutSecond) * time.Second,
			// TODO: User-Agent
		}
		defaultClient = MuseumClient{
			host:   host,
			Client: c,
		}
	})
	return defaultClient
}

type MuseumClient struct {
	host string
	http.Client

	// isPurge represents whether purge the `whole` tenant museum charts
	isPurge bool
	// specifiedVersion defines if we should operate the specified version
	specifiedVersion string
	// operation chart prefix
	prefix string
	// period represents the existing chart period
	period time.Duration
}

func (mc MuseumClient) Get(ctx context.Context, name string) ([]*Chart, error) {
	var cs []*Chart

	f, err := mc.get(ctx)
	if err != nil {
		return nil, err
	}

	for _, e := range f.Entries {
		if len(e) == 0 {
			continue
		}
		if e[0].Name != name {
			continue
		}
		for _, c := range e {
			cs = append(cs, NewChart(c))
		}
	}
	return cs, nil
}

func (mc MuseumClient) get(ctx context.Context) (*repo.IndexFile, error) {

	u, err := url.Parse(fmt.Sprintf("%s/index.yaml", mc.host))
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := mc.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f, err := hackLoadIndex(ctx, data)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (mc MuseumClient) GetAll(ctx context.Context) ([]*Chart, error) {

	var cs []*Chart

	f, err := mc.get(ctx)
	if err != nil {
		return nil, err
	}

	for _, e := range f.Entries {
		for _, c := range e {
			cs = append(cs, NewChart(c))
		}
	}
	return cs, nil
}

type DeleteOption func(MuseumClient) MuseumClient

func WithPurgeOption() DeleteOption {
	return func(c MuseumClient) MuseumClient {
		c.isPurge = true
		return c
	}
}

func WithPrefix(prefix string) DeleteOption {
	return func(c MuseumClient) MuseumClient {
		c.prefix = prefix
		return c
	}
}

func WithPeriod(period time.Duration) DeleteOption {
	return func(mc MuseumClient) MuseumClient {
		mc.period = period
		return mc
	}
}

// WithSpecifiedChartVersion set the specified chart version for the all chart operation
// e.g.: You can use it to fetch all charts matched with this provided version or delete them
func WithSpecifiedChartVersion(version string) DeleteOption {
	return func(c MuseumClient) MuseumClient {
		c.specifiedVersion = version
		return c
	}
}

func (mc MuseumClient) resolveAPIURL() (string, error) {
	u, err := url.Parse(mc.host)
	if err != nil {
		return "", err
	}
	u.Path = fmt.Sprintf("/api%scharts", u.Path)
	return u.String(), nil
}

func (mc MuseumClient) Del(ctx context.Context,
	charts []*Chart,
	opts ...DeleteOption,
) (int, error) {

	for _, opt := range opts {
		mc = opt(mc)
	}

	if mc.isPurge && len(charts) != 0 {
		return 0, ErrPurgeConflict
	}

	if len(charts) == 0 {
		cs, err := mc.GetAll(ctx)
		if err != nil {
			return 0, err
		}
		charts = append(charts, cs...)
	}
	var delCount int
	api, err := mc.resolveAPIURL()
	if err != nil {
		return 0, err
	}
	for _, chart := range charts {
		if mc.prefix != "" && !strings.HasPrefix(chart.Name, mc.prefix) {
			continue
		}
		version := chart.Version
		if mc.specifiedVersion != "" {
			version = mc.specifiedVersion
		}
		if mc.period != 0 && chart.Created.Unix() > time.Now().Add(mc.period).Unix() {
			continue
		}
		u, err := url.Parse(fmt.Sprintf("%s/%s/%s", api, chart.Name, version))
		if err != nil {
			return 0, err
		}
		req, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
		if err != nil {
			return 0, err
		}
		resp, err := mc.Do(req)
		if err != nil {
			log.Printf("WARNING: delete failed : Chart : %s:%s ,err: %q", chart.Name, chart.Version, err)
			continue
		}
		defer resp.Body.Close()
		log.Printf("DELETING CHART: %s:%s,Created: %s,STATUS: %d", chart.Name, chart.Version, chart.Created.String(), resp.StatusCode)
		delCount++
	}
	return delCount, nil
}
