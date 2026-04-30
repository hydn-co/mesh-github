package collectors

import (
	"context"
	"log/slog"
)

func logCollector(ctx context.Context, name string) {
	slog.InfoContext(ctx, "starting collector", "collector", name)
}
