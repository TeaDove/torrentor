package torrent_supplier

import (
	"context"
	stderrors "errors"
	"time"

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

func (r *Supplier) Stats(ctx context.Context, d time.Duration) <-chan torrent.ClientStats {
	ctx, cancel := context.WithTimeout(ctx, d)
	// TODO move to settings
	ticker := time.NewTicker(time.Second)
	statsChan := make(chan torrent.ClientStats)

	go func() {
		defer close(statsChan)
		defer cancel()
		// possible mem lick
		for {
			select {
			case <-ticker.C:
				statsChan <- r.client.Stats()
			case <-ctx.Done():
				return
			}
		}
	}()

	return statsChan
}
