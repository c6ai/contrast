// Copyright 2024 Edgeless Systems GmbH
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/edgelesssys/contrast/coordinator/internal/authority"
	"github.com/edgelesssys/contrast/internal/ca"
	"github.com/edgelesssys/contrast/internal/manifest"
	"github.com/edgelesssys/contrast/internal/memstore"
	"github.com/edgelesssys/contrast/internal/userapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type userAPIServer struct {
	policyTextStore store[manifest.HexString, manifest.Policy]
	manifSetGetter  manifestSetGetter
	logger          *slog.Logger

	userapi.UnimplementedUserAPIServer
}

func newUserAPIServer(mSGetter manifestSetGetter, log *slog.Logger) *userAPIServer {
	s := &userAPIServer{
		policyTextStore: memstore.New[manifest.HexString, manifest.Policy](),
		manifSetGetter:  mSGetter,
		logger:          log.WithGroup("userapi"),
	}

	return s
}

func (s *userAPIServer) SetManifest(ctx context.Context, req *userapi.SetManifestRequest,
) (*userapi.SetManifestResponse, error) {
	s.logger.Info("SetManifest called")

	if err := s.validatePeer(ctx); err != nil {
		s.logger.Warn("SetManifest peer validation failed", "err", err)
		return nil, status.Errorf(codes.PermissionDenied, "validating peer: %v", err)
	}

	var m *manifest.Manifest
	if err := json.Unmarshal(req.Manifest, &m); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unmarshaling manifest: %v", err)
	}

	if len(m.Policies) != len(req.Policies) {
		return nil, status.Error(codes.InvalidArgument, "request must contain exactly the policies referenced in the manifest")
	}

	for _, policyBytes := range req.Policies {
		policy := manifest.Policy(policyBytes)
		if _, ok := m.Policies[policy.Hash()]; !ok {
			return nil, status.Errorf(codes.InvalidArgument, "policy %v not found in manifest", policy.Hash())
		}
		s.policyTextStore.Set(policy.Hash(), policy)
	}

	ca, err := s.manifSetGetter.SetManifest(req.GetManifest(), req.GetPolicies())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "setting manifest: %v", err)
	}

	resp := &userapi.SetManifestResponse{
		RootCA: ca.GetRootCACert(),
		MeshCA: ca.GetMeshCACert(),
	}

	s.logger.Info("SetManifest succeeded")
	return resp, nil
}

func (s *userAPIServer) GetManifests(_ context.Context, _ *userapi.GetManifestsRequest,
) (*userapi.GetManifestsResponse, error) {
	s.logger.Info("GetManifest called")

	manifests, ca, err := s.manifSetGetter.GetManifestsAndLatestCA()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting manifests: %v", err)
	}
	if len(manifests) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "no manifests set")
	}

	manifestBytes, err := manifestSliceToBytesSlice(manifests)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshaling manifests: %v", err)
	}

	// TODO(burgerdev): these should be loaded from history.
	policies := s.policyTextStore.GetAll()
	if len(policies) == 0 {
		return nil, status.Error(codes.Internal, "no policies found in store")
	}

	resp := &userapi.GetManifestsResponse{
		Manifests: manifestBytes,
		Policies:  policySliceToBytesSlice(policies),
		RootCA:    ca.GetRootCACert(),
		MeshCA:    ca.GetMeshCACert(),
	}

	s.logger.Info("GetManifest succeeded")
	return resp, nil
}

func (s *userAPIServer) validatePeer(ctx context.Context) error {
	latest, err := s.manifSetGetter.LatestManifest()
	if err != nil && errors.Is(err, authority.ErrNoManifest) {
		// in the initial state, no peer validation is required
		return nil
	} else if err != nil && errors.Is(err, authority.ErrNeedsRecovery) {
		// TODO(burgerdev): give the user something more palatable.
		return err
	} else if err != nil {
		return fmt.Errorf("getting latest manifest: %w", err)
	}
	if len(latest.WorkloadOwnerKeyDigests) == 0 {
		return errors.New("setting manifest is disabled")
	}

	peerPubKey, err := getPeerPublicKey(ctx)
	if err != nil {
		return err
	}
	peerPub256Sum := sha256.Sum256(peerPubKey)
	for _, key := range latest.WorkloadOwnerKeyDigests {
		trustedWorkloadOwnerSHA256, err := key.Bytes()
		if err != nil {
			return fmt.Errorf("parsing key: %w", err)
		}
		if bytes.Equal(peerPub256Sum[:], trustedWorkloadOwnerSHA256) {
			return nil
		}
	}
	return errors.New("peer not authorized workload owner")
}

func getPeerPublicKey(ctx context.Context) ([]byte, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("no peer found in context")
	}
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, errors.New("peer auth info is not of type TLSInfo")
	}
	if len(tlsInfo.State.PeerCertificates) == 0 || tlsInfo.State.PeerCertificates[0] == nil {
		return nil, errors.New("no peer certificates found")
	}
	if tlsInfo.State.PeerCertificates[0].PublicKeyAlgorithm != x509.ECDSA {
		return nil, errors.New("peer public key is not of type ECDSA")
	}
	return x509.MarshalPKIXPublicKey(tlsInfo.State.PeerCertificates[0].PublicKey)
}

func policySliceToBytesSlice(s []manifest.Policy) [][]byte {
	var policies [][]byte
	for _, policy := range s {
		policies = append(policies, policy)
	}
	return policies
}

func manifestSliceToBytesSlice(s []*manifest.Manifest) ([][]byte, error) {
	var manifests [][]byte
	for i, manifest := range s {
		manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("mashaling manifest %d manifest: %w", i, err)
		}
		manifests = append(manifests, manifestBytes)
	}
	return manifests, nil
}

type manifestSetGetter interface {
	SetManifest(manifest []byte, policies [][]byte) (*ca.CA, error)
	GetManifestsAndLatestCA() ([]*manifest.Manifest, *ca.CA, error)
	LatestManifest() (*manifest.Manifest, error)
}

type store[keyT comparable, valueT any] interface {
	Get(key keyT) (valueT, bool)
	GetAll() []valueT
	Set(key keyT, value valueT)
}
