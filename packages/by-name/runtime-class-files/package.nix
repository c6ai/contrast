# Copyright 2024 Edgeless Systems GmbH
# SPDX-License-Identifier: AGPL-3.0-only

{ lib
, fetchurl
, stdenvNoCC
, igvmmeasure
, kata-image
, kata-igvm
, kata-runtime
, debugRuntime ? true
}:
let
  # Currently, those are files extracted from the CoCo AKS node image (AKSCBLMariner-V2katagen2).
  # In the future, those will be generated by us.
  rootfs = kata-image;
  igvm = if debugRuntime then kata-igvm.debug else kata-igvm;
  cloud-hypervisor-bin = fetchurl {
    url = "https://cdn.confidential.cloud/contrast/node-components/2024-03-13/cloud-hypervisor-cvm";
    hash = "sha256-coTHzd5/QLjlPQfrp9d2TJTIXKNuANTN7aNmpa8PRXo=";
  };
  containerd-shim-contrast-cc-v2 = "${kata-runtime}/bin/containerd-shim-kata-v2";
in
stdenvNoCC.mkDerivation {
  name = "runtime-class-files";
  version = "2024-03-13";

  dontUnpack = true;

  buildInputs = [ igvmmeasure ];

  buildPhase = ''
    mkdir -p $out
    igvmmeasure -b ${igvm} | dd conv=lcase > $out/launch-digest.hex
    echo -n "contrast-cc-" > $out/runtime-handler
    cat $out/launch-digest.hex | head -c 32 >> $out/runtime-handler
  '';

  passthru = {
    inherit debugRuntime rootfs igvm cloud-hypervisor-bin containerd-shim-contrast-cc-v2;
  };
}
