package main

import (
	"context"
	"log/slog"

	"github.com/edgelesssys/contrast/internal/recoveryapi"
)

type recoveryAPIServer struct {
	logger      *slog.Logger
	recoverable recoverable

	recoveryapi.UnimplementedRecoveryAPIServer
}

func newRecoveryAPIServer(recoveryTarget recoverable, log *slog.Logger) *recoveryAPIServer {
	s := &recoveryAPIServer{
		logger:      log.WithGroup("recoveryapi"),
		recoverable: recoveryTarget,
	}

	return s
}

func (s *recoveryAPIServer) Recover(_ context.Context, req *recoveryapi.RecoverRequest) (*recoveryapi.RecoverResponse, error) {
	s.logger.Info("Recover called")
	return &recoveryapi.RecoverResponse{}, s.recoverable.Recover(req.Seed, req.Salt)
}

type recoverable interface {
	Recover(seed, salt []byte) error
}
