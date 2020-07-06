package client

import (
	"context"

	"helm.sh/helm/pkg/repo"
	"sigs.k8s.io/yaml"
)

func hackLoadIndex(_ context.Context, index []byte) (*repo.IndexFile, error) {
	i := &repo.IndexFile{}
	if err := yaml.Unmarshal(index, i); err != nil {
		return i, err
	}
	i.SortEntries()
	return i, nil
}
