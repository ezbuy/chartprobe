package client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var ErrPurgeConflict = errors.New("client: purge only can be executed with empty charts list")

var ClientOnce sync.Once
var defaultClient MuseumClient

func NewClient() MuseumClient {
	ClientOnce.Do(func() {
		// TODO: multi-tenant
		host, ok := viper.Get("host").(string)
		if !ok {
			panic("client init: bad host")
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
}

func (mc MuseumClient) ensureIndex(ctx context.Context, chart *Chart) error {
	if chart == nil {
		return nil
	}
	return nil
}

func (mc MuseumClient) Get() {}
func (mc MuseumClient) GetAll(ctx context.Context) ([]*Chart, error) {

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

	var cs []*Chart

	for _, e := range f.Entries {
		for _, c := range e {
			cs = append(cs, NewChart(c.Name, c.Version))
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

func (mc MuseumClient) resolveAPIURL() (string, error) {
	u, err := url.Parse(mc.host)
	if err != nil {
		return "", err
	}
	u.Path = fmt.Sprintf("/api%scharts", u.Path)
	return u.String(), nil
}

func (mc MuseumClient) Del(ctx context.Context, charts []*Chart, opts ...DeleteOption) (int, error) {

	for _, opt := range opts {
		mc = opt(mc)
	}

	if mc.isPurge && len(charts) != 0 {
		return 0, ErrPurgeConflict
	}

	if mc.isPurge && len(charts) == 0 {
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
		u, err := url.Parse(fmt.Sprintf("%s/%s/%s", api, chart.Name, chart.Version))
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
		log.Printf("DELETING CHART: %s:%s,STATUS: %d", chart.Name, chart.Version, resp.StatusCode)
		delCount++
	}
	return len(charts), nil
}
