package assert

import (
	"log/slog"
	"os"
)

func NotError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

func True(condition bool) {

}
