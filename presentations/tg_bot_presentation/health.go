package tg_bot_presentation

import (
	"context"

	"github.com/pkg/errors"
)

func (r *Presentation) Health(_ context.Context) error {
	_, err := r.bot.GetMe()
	return errors.Wrap(err, "failed to get me")
}
