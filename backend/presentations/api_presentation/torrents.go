package api_presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (r *Presentation) statsTorrents(c *fiber.Ctx) error {
	serviceStats, _, err := r.torrentorService.Stats(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to get stats")
	}

	return c.JSON(fiber.Map{"serviceStats": serviceStats})
}

func (r *Presentation) download(c *fiber.Ctx) error {
	type Request struct {
		Magnet string `json:"magnet" validate:"required"`
	}

	req, err := parseJSON[Request](c)
	if err != nil {
		return errors.WithStack(err)
	}

	torrent, err := r.torrentorService.DownloadAndSaveFromMagnet(c.UserContext(), req.Magnet)
	if err != nil {
		return errors.Wrap(err, "failed to start download")
	}

	return c.JSON(torrent)
}

func (r *Presentation) listTorrents(c *fiber.Ctx) error {
	torrents, err := r.torrentorService.ListOpenTorrents(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to get torrents")
	}

	return c.JSON(torrents)
}

func (r *Presentation) getTorrent(c *fiber.Ctx) error {
	infoHash, err := getInfoHashParams(c)
	if err != nil {
		return errors.WithStack(err)
	}

	torrent, ok := r.torrentorService.GetTorrentByInfoHash(infoHash)
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "torrent not found"}
	}

	return c.JSON(torrent)
}

func (r *Presentation) getTorrentFile(c *fiber.Ctx) error {
	infoHash, err := getInfoHashParams(c)
	if err != nil {
		return errors.WithStack(err)
	}

	pathHash := c.Params("pathHash")

	torrent, ok := r.torrentorService.GetTorrentByInfoHash(infoHash)
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "torrent not found"}
	}

	file, ok := torrent.FileHashMap[pathHash]
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "file not found"}
	}

	return c.SendFile(file.Location())
}

func (r *Presentation) deleteTorrent(c *fiber.Ctx) error {
	infoHash, err := getInfoHashParams(c)
	if err != nil {
		return errors.WithStack(err)
	}

	ok := r.torrentorService.DeleteTorrentByInfoHash(infoHash)
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "torrent not found"}
	}

	return c.JSON(fiber.Map{"success": true})
}
