package bot

import "context"

func (tg *TelegramBotService) HelloCommand(ctx context.Context, args []string) error {
	panic("Hello I panicked!")
}