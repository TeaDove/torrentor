package web_app_presentation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func IndexForm(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{"IP": c.IP()})
}

func viewError(c fiber.Ctx, err error) error {
	return c.Render("error", fiber.Map{"Error": errors.Wrap(err, "failed to parse id")})
}

func (r *Presentation) TorrentForm(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse id"))
	}

	ctx := logger_utils.AddLoggerToCtx(c.Context())

	torrent, err := r.torrentorViewerService.GetTorrentMetadataByID(ctx, id)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get torrent metadata"))
	}

	return c.Render("torrent", fiber.Map{"TorrentName": torrent.Name})
}
