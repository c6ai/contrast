// Copyright 2024 Edgeless Systems GmbH
// SPDX-License-Identifier: AGPL-3.0-only

package authority

import (
	"context"

	"github.com/edgelesssys/contrast/coordinator/internal/seedengine"
	"github.com/edgelesssys/contrast/internal/recoveryapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recover recovers the Coordinator from a seed and salt.
func (a *Authority) Recover(_ context.Context, req *recoveryapi.RecoverRequest) (*recoveryapi.RecoverResponse, error) {
	a.logger.Info("Recover called")

	seedEngine, err := seedengine.New(req.Seed, req.Salt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "creating seed engine: %v", err)
	}
	if !a.se.CompareAndSwap(nil, seedEngine) {
		return nil, status.Error(codes.FailedPrecondition, "coordinator is already recovered")
	}
	a.hist.ConfigureSigningKey(a.se.Load().TransactionSigningKey())
	return &recoveryapi.RecoverResponse{}, nil
}
