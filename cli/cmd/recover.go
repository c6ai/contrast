package cmd

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/edgelesssys/contrast/internal/atls"
	"github.com/edgelesssys/contrast/internal/attestation/snp"
	"github.com/edgelesssys/contrast/internal/fsstore"
	"github.com/edgelesssys/contrast/internal/grpc/dialer"
	"github.com/edgelesssys/contrast/internal/recoveryapi"
	"github.com/edgelesssys/contrast/internal/userapi"
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
	cmd.Flags().String("workload-owner-key", workloadOwnerPEM, "path to workload owner key (.pem) file")
	cmd.Flags().String("seedshare-owner-key", seedshareOwnerPEM, "private key to decrypt the seed share")
	cmd.Flags().String("seed", seedSharesFilename, "file with the encrypted seed shares")

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

	workloadOwnerKey, err := loadWorkloadOwnerKey(flags.workloadOwnerKeyPath, nil, log)
	if err != nil {
		return fmt.Errorf("loading workload owner key: %w", err)
	}
	seed, salt, err := decryptedSeedFromShares(flags.seedSharesFilename, flags.seedShareOwnerKeyPath)
	if err != nil {
		return fmt.Errorf("decrypting seed: %w", err)
	}

	kdsDir, err := cachedir("kds")
	if err != nil {
		return fmt.Errorf("getting cache dir: %w", err)
	}
	log.Debug("Using KDS cache dir", "dir", kdsDir)

	validateOptsGen := newCoordinatorValidateOptsGen(flags.policy)
	kdsCache := fsstore.New(kdsDir, log.WithGroup("kds-cache"))
	kdsGetter := snp.NewCachedHTTPSGetter(kdsCache, snp.NeverGCTicker, log.WithGroup("kds-getter"))
	validator := snp.NewValidator(validateOptsGen, kdsGetter, log.WithGroup("snp-validator"))
	dialer := dialer.NewWithKey(atls.NoIssuer, validator, &net.Dialer{}, workloadOwnerKey)

	log.Debug("Dialing coordinator", "endpoint", flags.coordinator)
	conn, err := dialer.Dial(cmd.Context(), flags.coordinator)
	if err != nil {
		return fmt.Errorf("dialing coordinator: %w", err)
	}
	defer conn.Close()

	client := recoveryapi.NewRecoveryAPIClient(conn)
	req := &recoveryapi.RecoverRequest{
		Seed: seed,
		Salt: salt,
	}
	if _, err := client.Recover(cmd.Context(), req); err != nil {
		return fmt.Errorf("recovering: %w", err)
	}
	log.Debug("Got response")

	fmt.Fprintln(cmd.OutOrStdout(), "✔️ Successfully recovered the Coordinator")
	return nil
}

type recoverFlags struct {
	coordinator           string
	policy                []byte
	seedSharesFilename    string
	seedShareOwnerKeyPath string
	workloadOwnerKeyPath  string
}

func decryptedSeedFromShares(seedSharesPath, seedShareOwnerKeyPath string) ([]byte, []byte, error) {
	key, err := loadSeedshareOwnerKey(seedShareOwnerKeyPath)
	if err != nil {
		return nil, nil, err
	}
	pub, ok := key.Public().(*rsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("could not get public key from seedshare owner key")
	}
	var seedShareDoc userapi.SeedShareDocument
	seedShareBytes, err := os.ReadFile(seedSharesPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading seed shares: %w", err)
	}
	if err := json.Unmarshal(seedShareBytes, &seedShareDoc); err != nil {
		return nil, nil, fmt.Errorf("unmarshaling seed shares: %w", err)
	}
	for _, share := range seedShareDoc.SeedShares {
		shareKey, err := x509.ParsePKCS1PublicKey(share.PublicKey)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing seed share key: %w", err)
		}
		if !pub.Equal(shareKey) {
			continue
		}
		seed, err := rsa.DecryptOAEP(sha256.New(), nil, key, share.EncryptedSeed, []byte("seedshare"))
		if err != nil {
			return nil, nil, fmt.Errorf("decrypting seed share: %w", err)
		}
		return seed, seedShareDoc.Salt, nil
	}
	return nil, nil, fmt.Errorf("no matching seed share found")
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
	seed, err := cmd.Flags().GetString("seed")
	if err != nil {
		return recoverFlags{}, err
	}
	seedShareOwnerKeyPath, err := cmd.Flags().GetString("seedshare-owner-key")
	if err != nil {
		return recoverFlags{}, err
	}
	workloadOwnerKeyPath, err := cmd.Flags().GetString("workload-owner-key")
	if err != nil {
		return recoverFlags{}, err
	}

	return recoverFlags{
		coordinator:           coordinator,
		policy:                policy,
		seedSharesFilename:    seed,
		seedShareOwnerKeyPath: seedShareOwnerKeyPath,
		workloadOwnerKeyPath:  workloadOwnerKeyPath,
	}, nil
}

func loadSeedshareOwnerKey(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading seedshare owner key: %w", err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("decoding seedshare owner key: no key found")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("decoding seedshare owner key: invalid key type %q", block.Type)
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing seedshare owner key: %w", err)
	}
	return key, nil
}
