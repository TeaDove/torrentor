package torrent_supplier

import (
	"context"
	stderrors "errors"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
)

type Supplier struct {
	client *torrent.Client
}

func NewSupplier(_ context.Context, dataDir string) (*Supplier, error) {
	config := torrent.NewDefaultClientConfig()
	config.DataDir = dataDir

	client, err := torrent.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	return &Supplier{client: client}, nil
}

func (r *Supplier) Close(_ context.Context) error {
	errs := r.client.Close()
	if len(errs) != 0 {
		return errors.Wrap(stderrors.Join(errs...), "failed to close client")
	}

	return nil
}
