{ lib
, rustPlatform
, fetchFromGitHub
, cmake
, pkg-config
, protobuf
, withSeccomp ? true
, libseccomp
, withAgentPolicy ? false
, withStandardOCIRuntime ? false
}:
rustPlatform.buildRustPackage rec {
  pname = "kata-agent";
  version = "3.2.0";

  src = fetchFromGitHub {
    owner = "kata-containers";
    repo = "kata-containers";
    rev = version;
    hash = "sha256-zEKuEjw8a5hRNULNSkozjuHT6+hcbuTIbVPPImy/TsQ=";
  };

  sourceRoot = "${src.name}/src/agent";

  cargoHash = "sha256-m4Q3N1O5ME7V4I4c8tJtr/rGN4zpDe4p0c2s4mLeFuY=";

  nativeBuildInputs = [
    cmake
    pkg-config
    protobuf
  ];

  buildInputs = lib.optionals withAgentPolicy [
    libseccomp.dev
    libseccomp.lib
    libseccomp
  ];

  # Build.rs writes to src
  postConfigure = ''
    chmod -R +w ../..
  '';

  LIBC = "gnu";
  SECCOMP = if withSeccomp then "yes" else "no";
  AGENT_POLICY = if withAgentPolicy then "yes" else "no";
  STANDARD_OCI_RUNTIME = if withStandardOCIRuntime then "yes" else "no";

  buildPhase = ''
    runHook preBuild

    make

    runHook postBuild
  '';

  checkFlags = [
    "--skip=mount::tests::test_already_baremounted"
    "--skip=netlink::tests::list_routes stdout"
  ];
}
