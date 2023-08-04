package tools

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	quic "github.com/refraction-networking/uquic"
	"github.com/refraction-networking/uquic/internal/utils"
	"github.com/refraction-networking/uquic/logging"
	"github.com/refraction-networking/uquic/qlog"
)

func NewQlogger(logger io.Writer) func(context.Context, logging.Perspective, quic.ConnectionID) logging.ConnectionTracer {
	return func(_ context.Context, p logging.Perspective, connID quic.ConnectionID) logging.ConnectionTracer {
		role := "server"
		if p == logging.PerspectiveClient {
			role = "client"
		}
		filename := fmt.Sprintf("log_%x_%s.qlog", connID.Bytes(), role)
		fmt.Fprintf(logger, "Creating %s.\n", filename)
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to create qlog file: %s", err)
			return nil
		}
		bw := bufio.NewWriter(f)
		return qlog.NewConnectionTracer(utils.NewBufferedWriteCloser(bw, f), p, connID)
	}
}
