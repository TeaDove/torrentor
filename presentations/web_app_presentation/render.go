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

	torrent, err := r.torrentorService.GetTorrentMetadataByID(c.Context(), id)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get torrent metadata"))
	}

	return c.Render("torrent", fiber.Map{"TorrentName": torrent.Name, "TorrentFiles": torrent.Root.FlatFiles()})
}

func (r *Presentation) FileForm(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse id"))
	}

	file, err := r.torrentorService.GetFileContentByID(c.Context(), id)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get file content"))
	}

	fileStats, err := file.Stat()
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get file stats"))
	}

	if fileStats.IsDir() {
		return viewError(c, errors.New("file is a directory"))
	}

	return c.SendStream(file, int(fileStats.Size()))
}
