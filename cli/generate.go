package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/katexochen/coordinator-kbs/internal/kubeapi"
	"github.com/katexochen/coordinator-kbs/internal/manifest"
	"github.com/spf13/cobra"
)

const kataPolicyAnnotationKey = "io.katacontainers.config.agent.policy"

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate",
		RunE:  runGenerate,
	}

	cmd.Flags().StringP("policy", "p", "", "path to policy (.rego) file")
	cobra.MarkFlagRequired(cmd.Flags(), "policy")
	cmd.Flags().StringP("settings", "s", "", "path to settings (.json) file")
	cobra.MarkFlagRequired(cmd.Flags(), "settings")
	cmd.Flags().StringP("manifest", "m", "", "path to manifest (.json) file")
	cobra.MarkFlagRequired(cmd.Flags(), "manifest")

	return cmd
}

func runGenerate(cmd *cobra.Command, args []string) error {
	flags, err := parseGenerateFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	paths, err := findGenerateTargets(args)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no .yml/.yaml files found")
	}

	if err := generatePolicies(cmd.Context(), flags.policyPath, flags.settingsPath, paths); err != nil {
		return fmt.Errorf("failed to generate policies: %w", err)
	}

	policies, err := manifestPoliciesFromKubeResources(paths)
	if err != nil {
		return fmt.Errorf("failed to find kube resources with policy: %w", err)
	}

	manifestData, err := os.ReadFile(flags.manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}
	var manifest *manifest.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to unmarshal manifest: %w", err)
	}
	manifest.Policies = policies
	manifestData, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	if err := os.WriteFile(flags.manifestPath, manifestData, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

func findGenerateTargets(args []string) ([]string, error) {
	var paths []string
	for _, path := range args {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil // Skip directories
			}
			switch {
			case strings.HasSuffix(info.Name(), ".yaml"):
				paths = append(paths, path)
			case strings.HasSuffix(info.Name(), ".yml"):
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk %s: %w", path, err)
		}
	}
	return paths, nil
}

func generatePolicies(ctx context.Context, regoPath, policyPath string, yamlPaths []string) error {
	for _, yamlPath := range yamlPaths {
		policyHash, err := generatePolicyForFile(ctx, regoPath, policyPath, yamlPath)
		if err != nil {
			return fmt.Errorf("failed to generate policy for %s: %w", yamlPath, err)
		}
		if policyHash == [32]byte{} {
			continue
		}
		fmt.Printf("%x  %s\n", policyHash, yamlPath)
	}
	return nil
}

func generatePolicyForFile(ctx context.Context, regoPath, policyPath, yamlPath string) ([32]byte, error) {
	args := []string{
		"--raw-out",
		"--use-cached-files",
		fmt.Sprintf("--input-files-path=%s", regoPath),
		fmt.Sprintf("--settings-file-name=%s", policyPath),
		fmt.Sprintf("--yaml-file=%s", yamlPath),
	}
	genpolicy := exec.CommandContext(ctx, genpolicyPath, args...)
	var stdout, stderr bytes.Buffer
	genpolicy.Stdout = &stdout
	genpolicy.Stderr = &stderr
	if err := genpolicy.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return [32]byte{}, fmt.Errorf("genpolicy failed with exit code %d: %s",
				exitErr.ExitCode(), stderr.String())
		}
		return [32]byte{}, fmt.Errorf("genpolicy failed: %w", err)
	}
	if stdout.Len() == 0 {
		log.Printf("policy output for %s is empty, ignoring the file", yamlPath)
		return [32]byte{}, nil
	}
	policyHash := sha256.Sum256(stdout.Bytes())
	return policyHash, nil
}

func manifestPoliciesFromKubeResources(yamlPaths []string) (map[manifest.HexString]string, error) {
	var kubeObjs []any
	for _, path := range yamlPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
		objs, err := kubeapi.UnmarshalK8SResources(data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", path, err)
		}
		kubeObjs = append(kubeObjs, objs...)
	}

	policies := make(map[manifest.HexString]string)
	for _, objAny := range kubeObjs {
		var name, annotation string
		switch obj := objAny.(type) {
		case kubeapi.Pod:
			name = obj.Name
			annotation = obj.Annotations[kataPolicyAnnotationKey]
		case kubeapi.Deployment:
			name = obj.Name
			annotation = obj.Spec.Template.Annotations[kataPolicyAnnotationKey]
		case kubeapi.ReplicaSet:
			name = obj.Name
			annotation = obj.Spec.Template.Annotations[kataPolicyAnnotationKey]
		case kubeapi.StatefulSet:
			name = obj.Name
			annotation = obj.Spec.Template.Annotations[kataPolicyAnnotationKey]
		case kubeapi.DaemonSet:
			name = obj.Name
			annotation = obj.Spec.Template.Annotations[kataPolicyAnnotationKey]
		}
		if annotation == "" {
			continue
		}
		if name == "" {
			return nil, fmt.Errorf("name is required but empty")
		}
		policy, err := base64.StdEncoding.DecodeString(annotation)
		if err != nil {
			return nil, fmt.Errorf("failed to decode policy for %s: %w", name, err)
		}
		policyHash := sha256.Sum256(policy)
		policyHashStr := manifest.NewHexString(policyHash[:])
		if existingName, ok := policies[policyHashStr]; ok {
			if existingName != name {
				return nil, fmt.Errorf("policy hash collision: %s and %s have the same hash %s",
					existingName, name, policyHashStr)
			}
			continue
		}
		policies[policyHashStr] = name
		fmt.Printf("%s  %s\n", policyHashStr, name)
	}

	return policies, nil
}

type generateFlags struct {
	policyPath   string
	settingsPath string
	manifestPath string
}

func parseGenerateFlags(cmd *cobra.Command) (*generateFlags, error) {
	policyPath, err := cmd.Flags().GetString("policy")
	if err != nil {
		return nil, err
	}
	settingsPath, err := cmd.Flags().GetString("settings")
	if err != nil {
		return nil, err
	}
	manifestPath, err := cmd.Flags().GetString("manifest")
	if err != nil {
		return nil, err
	}
	return &generateFlags{
		policyPath:   policyPath,
		settingsPath: settingsPath,
		manifestPath: manifestPath,
	}, nil
}