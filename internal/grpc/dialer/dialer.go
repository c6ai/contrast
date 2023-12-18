/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

// Package dialer provides a grpc dialer that can be used to create grpc client connections with different levels of ATLS encryption / verification.
package dialer

import (
	"context"
	"crypto/ecdsa"
	"net"
	"time"

	"github.com/edgelesssys/nunki/internal/atls"
	"github.com/edgelesssys/nunki/internal/grpc/atlscredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dialer can open grpc client connections with different levels of ATLS encryption / verification.
type Dialer struct {
	issuer    atls.Issuer
	validator atls.Validator
	netDialer NetDialer
	privKey   *ecdsa.PrivateKey
}

// New creates a new Dialer.
func New(issuer atls.Issuer, validator atls.Validator, netDialer NetDialer) *Dialer {
	return &Dialer{
		issuer:    issuer,
		validator: validator,
		netDialer: netDialer,
	}
}

func NewWithKey(issuer atls.Issuer, validator atls.Validator, netDialer NetDialer, privKey *ecdsa.PrivateKey) *Dialer {
	return &Dialer{
		issuer:    issuer,
		validator: validator,
		netDialer: netDialer,
		privKey:   privKey,
	}
}

// Dial creates a new grpc client connection to the given target using the atls validator.
func (d *Dialer) Dial(ctx context.Context, target string) (*grpc.ClientConn, error) {
	var validators []atls.Validator
	if d.validator != nil {
		validators = append(validators, d.validator)
	}
	credentials := atlscredentials.NewWithKey(d.issuer, validators, d.privKey)

	return grpc.DialContext(ctx, target,
		d.grpcWithDialer(),
		grpc.WithTransportCredentials(credentials),
		grpc.WithConnectParams(grpc.ConnectParams{
			// We need a high initial timeout, because otherwise the client will get stuck in a reconnect loop
			// where the timeout is too low to get a full handshake done.
			MinConnectTimeout: 30 * time.Second,
		}),
	)
}

// DialInsecure creates a new grpc client connection to the given target without using encryption or verification.
// Only use this method when using another kind of encryption / verification (VPN, etc).
func (d *Dialer) DialInsecure(ctx context.Context, target string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target,
		d.grpcWithDialer(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

// DialNoVerify creates a new grpc client connection to the given target without verifying the server's attestation.
func (d *Dialer) DialNoVerify(ctx context.Context, target string) (*grpc.ClientConn, error) {
	credentials := atlscredentials.New(nil, nil)

	return grpc.DialContext(ctx, target,
		d.grpcWithDialer(),
		grpc.WithTransportCredentials(credentials),
	)
}

func (d *Dialer) grpcWithDialer() grpc.DialOption {
	return grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		return d.netDialer.DialContext(ctx, "tcp", addr)
	})
}

// NetDialer implements the net Dialer interface.
type NetDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}
