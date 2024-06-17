// Copyright 2024 Edgeless Systems GmbH
// SPDX-License-Identifier: AGPL-3.0-only

package authority

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/edgelesssys/contrast/internal/crypto"
	"github.com/edgelesssys/contrast/internal/manifest"
	"github.com/edgelesssys/contrast/internal/recoveryapi"
	"github.com/edgelesssys/contrast/internal/userapi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestRecovery(t *testing.T) {
	require := require.New(t)

	// Recover with an empty state should succeed.
	a1, _ := newAuthority(t)

	seedShareOwnerKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)
	seedShareOwnerKeyBytes := crypto.MarshalSeedShareOwnerKey(&seedShareOwnerKey.PublicKey)

	policy := []byte("=== SOME REGO HERE ===")
	policyHash := sha256.Sum256(policy)
	policyHashHex := manifest.NewHexString(policyHash[:])
	mnfst := &manifest.Manifest{
		Policies:              map[manifest.HexString][]string{policyHashHex: {"test"}},
		SeedshareOwnerPubKeys: []manifest.HexString{seedShareOwnerKeyBytes},
	}
	manifestBytes, err := json.Marshal(mnfst)
	require.NoError(err)

	req := &userapi.SetManifestRequest{
		Manifest: manifestBytes,
		Policies: [][]byte{policy},
	}
	resp1, err := a1.SetManifest(context.Background(), req)
	require.NoError(err)
	require.NotNil(resp1)
	seedSharesDoc := resp1.GetSeedSharesDoc()
	require.NotNil(seedSharesDoc)
	seedShares := seedSharesDoc.GetSeedShares()
	require.Len(seedShares, 1)

	seed, err := crypto.DecryptSeedShare(seedShareOwnerKey, seedShares[0])
	require.NoError(err)

	// A new authority on existing state should refuse all manifest calls.
	a2 := New(a1.hist, prometheus.NewRegistry(), slog.Default())
	_, err = a2.SetManifest(context.Background(), req)
	require.ErrorContains(err, ErrNeedsRecovery.Error())
	_, _, err = a2.getManifestsAndLatestCA()
	require.ErrorIs(err, ErrNeedsRecovery)

	// Recover with non-empty state should succeed.
	recoverReq := &recoveryapi.RecoverRequest{
		Seed: seed,
		Salt: seedSharesDoc.GetSalt(),
	}
	_, err = a2.Recover(context.Background(), recoverReq)
	require.NoError(err)
	manifests, _, err := a2.getManifestsAndLatestCA()
	require.NoError(err)
	require.Equal([]*manifest.Manifest{mnfst}, manifests)

	// Recover on a recovered authority should fail.
	_, err = a2.Recover(context.Background(), recoverReq)
	require.Error(err)
}
