// Copyright 2024 Edgeless Systems GmbH
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/edgelesssys/contrast/internal/ca"
	"github.com/edgelesssys/contrast/internal/crypto"
	"github.com/edgelesssys/contrast/internal/logger"
	"github.com/edgelesssys/contrast/internal/meshapi"
	"github.com/edgelesssys/contrast/internal/seedengine"
	"github.com/edgelesssys/contrast/internal/userapi"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (retErr error) {
	logger, err := logger.Default()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: creating logger: %v\n", err)
		return err
	}
	defer func() {
		if retErr != nil {
			logger.Error("Coordinator terminated after failure", "err", retErr)
		}
	}()

	logger.Info("Coordinator started")

	seed, err := crypto.GenerateRandomBytes(256)
	if err != nil {
		return fmt.Errorf("generating random seed: %w", err)
	}
	salt, err := crypto.GenerateRandomBytes(256)
	if err != nil {
		return fmt.Errorf("generating random salt: %w", err)
	}
	logger.Info("Generated new seed and salt")

	seedEngine, err := seedengine.New(seed, salt)
	if err != nil {
		return fmt.Errorf("creating seed engine: %w", err)
	}

	caInstance, err := ca.New(seedEngine.RootCAKey())
	if err != nil {
		return fmt.Errorf("creating CA: %w", err)
	}

	meshAuth := newMeshAuthority(caInstance, logger)

	eg := errgroup.Group{}

	eg.Go(func() error {
		userAPI := newUserAPIServer(meshAuth, caInstance, logger)
		logger.Info("Coordinator user API listening")
		if err := userAPI.Serve(net.JoinHostPort("0.0.0.0", userapi.Port)); err != nil {
			return fmt.Errorf("serving Coordinator API: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		meshAPI := newMeshAPIServer(meshAuth, caInstance, logger)
		logger.Info("Coordinator mesh API listening")
		if err := meshAPI.Serve(net.JoinHostPort("0.0.0.0", meshapi.Port)); err != nil {
			return fmt.Errorf("serving mesh API: %w", err)
		}
		return nil
	})

	return eg.Wait()
}
