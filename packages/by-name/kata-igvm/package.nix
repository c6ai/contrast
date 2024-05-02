{ lib
, fetchFromGitHub
, stdenv
, igvm-tooling
, kata-kernel-uvm
, kata-image
}:
let
  # This is not a real version, since the igvm builder is not part of the official release
  kata-version = "3.2.0.igvm";
  # keep up to date with the igvm-builder branch
  # https://github.com/microsoft/kata-containers/tree/dadelan/igvm-builder
  source = fetchFromGitHub {
    owner = "microsoft";
    repo = "kata-containers";
    rev = "ad93335ff0d1502a6f094324aa87275c8201c684";
    hash = "sha256-pVogv30WsQejBtheGz76O4MDUs1+nxm8Xr6LXGmtolg=";
  };
in
stdenv.mkDerivation rec {
  pname = "kata-containers-igvm";
  version = kata-version;

  nativeBuildInputs = [
    igvm-tooling
  ];

  src = source;

  sourceRoot = "${src.name}/tools/osbuilder/igvm-builder";

  postPatch = ''
    chmod +x igvm_builder.sh
    substituteInPlace igvm_builder.sh \
      --replace-fail '#!/usr/bin/env bash' '#!${stdenv.shell}' \
      --replace-fail 'python3 igvm/igvmgen.py' igvmgen \
      --replace-fail igvm/acpi/acpi-clh/ "${igvm-tooling}/share/igvm-tooling/acpi/acpi-clh/" \
      --replace-fail 'mv ''${igvm_name} ''${script_dir}' "" \
      --replace-fail sudo ""
    # TODO: cleanup
    #  --replace-fail '-acpi igvm/acpi/acpi-clh/' "" \
    # prevent non-hermetic download of igvm-tooling / igvmgen
    mkdir -p msigvm-1.2.0
  '';

  buildPhase = ''
    runHook preBuild

    # TODO: check if signature is deterministic
    ./igvm_builder.sh -k ${kata-kernel-uvm}/bzImage -v ${kata-image}/dm_verity.txt -o $out

    runHook postBuild
  '';
}
