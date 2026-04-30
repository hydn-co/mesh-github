package actions

import (
	"context"
	"log/slog"
)

func logAction(ctx context.Context, name string) {
	slog.InfoContext(ctx, "starting action", "action", name)
}
