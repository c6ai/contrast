# Copyright 2024 Edgeless Systems GmbH
# SPDX-License-Identifier: AGPL-3.0-only

{ buildGoModule
, fetchFromGitHub
, yq-go
, git
}:
buildGoModule rec {
  pname = "kata-runtime";
  version = "3.2.0.azl1";

  src = fetchFromGitHub {
    owner = "microsoft";
    repo = "kata-containers";
    rev = version;
    hash = "sha256-SsgI6h4/rjVWqUySoNgnbqAS9TdFAl05Fk2M1mJP3wM=";
  };

  sourceRoot = "${src.name}/src/runtime";

  preBuild = ''
    substituteInPlace Makefile \
      --replace-fail 'include golang.mk' ""
    for f in $(find . -name '*.in' -type f); do
      make ''${f%.in}
    done
  '';

  checkFlags =
    let
      # Skip tests that require a working hypervisor
      skippedTests = [
        "TestArchRequiredKernelModules"
        "TestCheckCLIFunctionFail"
        "TestEnvCLIFunction(|Fail)"
        "TestEnvGetAgentInfo"
        "TestEnvGetEnvInfo(|SetsCPUType|NoHypervisorVersion|AgentError|NoOSRelease|NoProcCPUInfo|NoProcVersion)"
        "TestEnvGetRuntimeInfo"
        "TestEnvHandleSettings(|InvalidParams)"
        "TestGetHypervisorInfo"
        "TestGetHypervisorInfoSocket"
        "TestSetCPUtype"
      ];
    in
    [ "-skip=^${builtins.concatStringsSep "$|^" skippedTests}$" ];

  CGO_ENABLED = 0;
  ldflags = [ "-s" ];

  vendorHash = null;
  subPackages = [
    "cmd/containerd-shim-kata-v2"
    "cmd/kata-monitor"
    # TODO(malt3): enable kata-runtime
    # It depends on CGO and kvm
    # "cmd/kata-runtime"
  ];

  nativeBuildInputs = [
    yq-go
    git
  ];

  meta.mainProgram = "containerd-shim-kata-v2";
}
