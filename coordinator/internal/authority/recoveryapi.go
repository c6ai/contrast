// Copyright 2024 Edgeless Systems GmbH
// SPDX-License-Identifier: AGPL-3.0-only

package authority

import (
	"context"
	"errors"
	"fmt"

	"github.com/edgelesssys/contrast/coordinator/internal/seedengine"
	"github.com/edgelesssys/contrast/internal/recoveryapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrAlreadyRecovered is returned if recovery was requested but a seed is already set.
var ErrAlreadyRecovered = errors.New("coordinator is already recovered")

// Recover recovers the Coordinator from a seed and salt.
func (a *Authority) Recover(_ context.Context, req *recoveryapi.RecoverRequest) (*recoveryapi.RecoverResponse, error) {
	a.logger.Info("Recover called")

	err := a.recover(req.Seed, req.Salt)
	switch {
	case errors.Is(err, ErrAlreadyRecovered):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case err == nil:
		return &recoveryapi.RecoverResponse{}, nil
	default:
		// Pretty sure this failed because the seed was bad.
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}
}

// recover recovers the seed engine from a seed and salt.
func (a *Authority) recover(seed, salt []byte) error {
	seedEngine, err := seedengine.New(seed, salt)
	if err != nil {
		return fmt.Errorf("creating seed engine: %w", err)
	}
	if !a.se.CompareAndSwap(nil, seedEngine) {
		return ErrAlreadyRecovered
	}
	a.hist.ConfigureSigningKey(seedEngine.TransactionSigningKey())
	return nil
}

func (a *Authority) checkSeedEngine() error {
	if a.se.Load() != nil {
		return nil
	}
	hasLatest, err := a.hist.HasLatest()
	if err != nil {
		return fmt.Errorf("querying history backend: %w", err)
	}
	if hasLatest {
		return ErrNeedsRecovery
	}
	return nil
}
