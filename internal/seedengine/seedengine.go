// Package seedengine provides deterministic key derivation of ECDSA and symmetric keys
// from a secret seed.
package seedengine

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"

	"filippo.io/keygen"
	"golang.org/x/crypto/hkdf"
)

// SeedEngine provides deterministic key derivation of ECDSA and symmetric keys
// from a secret seed.
type SeedEngine struct {
	curve   func() elliptic.Curve
	hashFun func() hash.Hash
	salt    []byte

	podStateSeed []byte
	historySeed  []byte

	rootCAKey             *ecdsa.PrivateKey
	transactionSigningKey *ecdsa.PrivateKey
}

// New creates a new SeedEngine from a secret seed and a salt.
func New(secretSeed []byte, salt []byte) (*SeedEngine, error) {
	se := &SeedEngine{
		curve:   elliptic.P256,
		hashFun: sha256.New,
		salt:    salt,
	}

	// Recommended to use salt length equal to hash size, see RFC 5869, section 3.1.
	if len(salt) != se.hashFun().Size() {
		return nil, fmt.Errorf("salt must be %d bytes long", se.hashFun().Size())
	}

	var err error
	se.podStateSeed, err = se.hkdfDerive(secretSeed, "POD STATE SECRET")
	if err != nil {
		return nil, fmt.Errorf("deriving seed: %w", err)
	}
	se.historySeed, err = se.hkdfDerive(secretSeed, "HISTORY SECRET")
	if err != nil {
		return nil, fmt.Errorf("deriving seed: %w", err)
	}
	transactionSigningSeed, err := se.hkdfDerive(secretSeed, "TRANSACTION SIGNING SECRET")
	if err != nil {
		return nil, fmt.Errorf("deriving seed: %w", err)
	}
	rootCASeed, err := se.hkdfDerive(secretSeed, "ROOT CA SEED")
	if err != nil {
		return nil, fmt.Errorf("deriving seed: %w", err)
	}

	se.transactionSigningKey, err = se.generateECDSAPrivateKey(transactionSigningSeed)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key: %w", err)
	}
	se.rootCAKey, err = se.generateECDSAPrivateKey(rootCASeed)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key: %w", err)
	}

	return se, nil
}

// DerivePodSecret derives a secret for a pod from the policy hash and the secret seed.
func (s *SeedEngine) DerivePodSecret(policyHash []byte) ([]byte, error) {
	return s.hkdfDerive(s.podStateSeed, fmt.Sprintf("POD SECRET %x", policyHash))
}

// DeriveMeshCAKey derives a secret for a mesh CA from the transaction and the secret seed.
func (s *SeedEngine) DeriveMeshCAKey(transaction []byte) (*ecdsa.PrivateKey, error) {
	transactionSecret, err := s.hkdfDerive(s.historySeed, fmt.Sprintf("TRANSACTION SECRET %x", transaction))
	if err != nil {
		return nil, err
	}
	meshCASeed, err := s.hkdfDerive(transactionSecret, "MESH CA SECRET")
	if err != nil {
		return nil, err
	}
	return s.generateECDSAPrivateKey(meshCASeed)
}

// RootCAKey returns the root CA key which is derived from the secret seed.
func (s *SeedEngine) RootCAKey() *ecdsa.PrivateKey {
	return s.rootCAKey
}

// TransactionSigningKey returns the transaction signing key which is derived from the secret seed.
func (s *SeedEngine) TransactionSigningKey() *ecdsa.PrivateKey {
	return s.transactionSigningKey
}

func (s *SeedEngine) hkdfDerive(secret []byte, info string) ([]byte, error) {
	hkdf := hkdf.New(s.hashFun, secret, s.salt, []byte(info))
	newSecret := make([]byte, len(secret))
	if _, err := io.ReadFull(hkdf, newSecret); err != nil {
		return nil, err
	}
	return newSecret, nil
}

func (s *SeedEngine) generateECDSAPrivateKey(secret []byte) (*ecdsa.PrivateKey, error) {
	return keygen.ECDSA(s.curve(), secret)
}
