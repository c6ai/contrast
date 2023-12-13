{ rustPlatform
, fetchFromGitHub
, cmake
, pkg-config
, protobuf
}:
rustPlatform.buildRustPackage rec {
  pname = "kata-runtime";
  version = "3.2.0";

  src = fetchFromGitHub {
    owner = "kata-containers";
    repo = "kata-containers";
    rev = version;
    hash = "sha256-zEKuEjw8a5hRNULNSkozjuHT6+hcbuTIbVPPImy/TsQ=";
  };

  sourceRoot = "${src.name}/src/runtime";

  cargoHash = "sha256-m4Q3N1O5ME7V4I4c8tJtr/rGN4zpDe4p0c2s4mLeFuY=";

  nativeBuildInputs = [
    cmake
    pkg-config
    protobuf
  ];

  # Build.rs writes to src
  postConfigure = ''
    chmod -R +w ../..
  '';

  LIBC = "gnu";
  SECCOMP = "no";

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
