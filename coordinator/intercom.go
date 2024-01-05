package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/edgelesssys/nunki/internal/atls"
	"github.com/edgelesssys/nunki/internal/attestation/snp"
	"github.com/edgelesssys/nunki/internal/grpc/atlscredentials"
	"github.com/edgelesssys/nunki/internal/intercom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"k8s.io/utils/clock"
)

type intercomServer struct {
	grpc          *grpc.Server
	certGet       certGetter
	caChainGetter certChainGetter
	ticker        clock.Ticker
	logger        *slog.Logger

	intercom.UnimplementedIntercomServer
}

type certGetter interface {
	GetCert(peerPublicKeyHashStr string) ([]byte, error)
}

func newIntercomServer(meshAuth *meshAuthority, caGetter certChainGetter, log *slog.Logger) *intercomServer {
	ticker := clock.RealClock{}.NewTicker(24 * time.Hour)
	validator := snp.NewValidatorWithCallbacks(meshAuth, ticker, log, meshAuth)
	credentials := atlscredentials.New(atls.NoIssuer, []atls.Validator{validator})
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials),
		grpc.KeepaliveParams(keepalive.ServerParameters{Time: 15 * time.Second}),
	)
	s := &intercomServer{
		grpc:          grpcServer,
		certGet:       meshAuth,
		caChainGetter: caGetter,
		ticker:        ticker,
		logger:        log.WithGroup("intercom"),
	}
	intercom.RegisterIntercomServer(s.grpc, s)
	return s
}

func (i *intercomServer) Serve(endpoint string) error {
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	defer i.ticker.Stop()
	return i.grpc.Serve(lis)
}

func (i *intercomServer) NewMeshCert(_ context.Context, req *intercom.NewMeshCertRequest,
) (*intercom.NewMeshCertResponse, error) {
	i.logger.Info("NewMeshCert called")

	cert, err := i.certGet.GetCert(req.PeerPublicKeyHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"getting certificate with public key hash %q: %v", req.PeerPublicKeyHash, err)
	}

	meshCACert := i.caChainGetter.GetMeshCACert()
	intermCert := i.caChainGetter.GetIntermCert()

	return &intercom.NewMeshCertResponse{
		MeshCACert: meshCACert,
		CertChain:  append(cert, intermCert...),
		RootCACert: i.caChainGetter.GetRootCACert(),
	}, nil
}
