package cmd

import (
	"fmt"
	"net"

	"github.com/edgelesssys/contrast/internal/atls"
	"github.com/edgelesssys/contrast/internal/attestation/snp"
	"github.com/edgelesssys/contrast/internal/fsstore"
	"github.com/edgelesssys/contrast/internal/grpc/dialer"
	"github.com/edgelesssys/contrast/internal/recoveryapi"
	"github.com/spf13/cobra"
)

// NewRecoverCmd creates the contrast recover subcommand.
func NewRecoverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover [flags]",
		Short: "recover a contrast deployment after restart",
		Long: `Recover a contrast deployment after restart.

The state of the Coordinator is stored protected on a persistent volume.
After a restart, the Coordinator requires the seed to derive the signing
key and verify the state integrity.

The recover command is used to provide the seed to the Coordinator.`,
		RunE: withTelemetry(runRecover),
	}

	cmd.Flags().StringP("coordinator", "c", "", "endpoint the coordinator can be reached at")
	must(cobra.MarkFlagRequired(cmd.Flags(), "coordinator"))
	cmd.Flags().String("coordinator-policy-hash", DefaultCoordinatorPolicyHash, "override the expected policy hash of the coordinator")

	return cmd
}

func runRecover(cmd *cobra.Command, _ []string) error {
	flags, err := parseRecoverFlags(cmd)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	log, err := newCLILogger(cmd)
	if err != nil {
		return err
	}
	log.Debug("Starting recovery")

	kdsDir, err := cachedir("kds")
	if err != nil {
		return fmt.Errorf("getting cache dir: %w", err)
	}
	log.Debug("Using KDS cache dir", "dir", kdsDir)

	validateOptsGen := newCoordinatorValidateOptsGen(flags.policy)
	kdsCache := fsstore.New(kdsDir, log.WithGroup("kds-cache"))
	kdsGetter := snp.NewCachedHTTPSGetter(kdsCache, snp.NeverGCTicker, log.WithGroup("kds-getter"))
	validator := snp.NewValidator(validateOptsGen, kdsGetter, log.WithGroup("snp-validator"))
	dialer := dialer.New(atls.NoIssuer, validator, &net.Dialer{})

	log.Debug("Dialing coordinator", "endpoint", flags.coordinator)
	conn, err := dialer.Dial(cmd.Context(), flags.coordinator)
	if err != nil {
		return fmt.Errorf("Error: failed to dial coordinator: %w", err)
	}
	defer conn.Close()

	client := recoveryapi.NewRecoveryAPIClient(conn)
	req := &recoveryapi.RecoverRequest{
		// TODO
	}
	if _, err := client.Recover(cmd.Context(), req); err != nil {
		return fmt.Errorf("recovering: %w", err)
	}
	log.Debug("Got response")

	fmt.Fprintln(cmd.OutOrStdout(), "✔️ Successfully recovered the Coordinator")
	return nil
}

type recoverFlags struct {
	coordinator string
	policy      []byte
}

func parseRecoverFlags(cmd *cobra.Command) (recoverFlags, error) {
	coordinator, err := cmd.Flags().GetString("coordinator")
	if err != nil {
		return recoverFlags{}, err
	}
	policy, err := decodeCoordinatorPolicyHash(cmd.Flags())
	if err != nil {
		return recoverFlags{}, err
	}

	return recoverFlags{
		coordinator: coordinator,
		policy:      policy,
	}, nil
}
