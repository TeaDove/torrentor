package tg_bot_presentation

import "context"

func (r *Presentation) Health(ctx context.Context) error {
	_, err := r.bot.GetMe()
	return err
}
