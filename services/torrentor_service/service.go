package torrentor_service

import (
	"context"
	"torrentor/repositories/torrent_repository"
	"torrentor/suppliers/torrent_supplier"
)

type Service struct {
	torrentSupplier   *torrent_supplier.Supplier
	torrentRepository *torrent_repository.Repository
}

func NewService(
	_ context.Context,
	torrentSupplier *torrent_supplier.Supplier,
	torrentRepository *torrent_repository.Repository,
) (*Service, error) {
	return &Service{
		torrentSupplier:   torrentSupplier,
		torrentRepository: torrentRepository,
	}, nil
}
