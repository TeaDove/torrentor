package web_app_presentation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func viewError(c fiber.Ctx, err error) error {
	zerolog.Ctx(c.Context()).
		Error().
		Stack().Err(err).
		Msg("view.error")

	return c.Render("error", fiber.Map{"Error": errors.Wrap(err, "failed to parse id")})
}
func IndexForm(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{"IP": c.IP()})
}

func (r *Presentation) TorrentForm(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse id"))
	}

	torrent, err := r.torrentorViewerService.GetTorrentMetadataByID(c.Context(), id)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get torrent metadata"))
	}

	return c.Render("torrent", fiber.Map{"TorrentName": torrent.Name, "TorrentFiles": torrent.Root.FlatFiles()})
}
