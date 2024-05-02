{ lib
, fetchurl
, linuxManualConfig
, fetchFromGitHub
, symlinkJoin
, stdenv
, patchutils
}:
let
  kver = "6.1.0";
  modDirVersion = "${kver}.mshv16";
  kataVersion = "3.2.0.azl1";
  tarfs_make = builtins.path {
    path = ./src;
  };
  tarfs_patch = fetchurl {
    name = "tarfs.patch";
    # update whenever tarfs.c changes: https://github.com/microsoft/kata-containers/commits/msft-main/src/tarfs/tarfs.c
    url = "https://raw.githubusercontent.com/microsoft/kata-containers/14e22b6d577425c2edc1318d77dd44828e9fec62/src/tarfs/tarfs.c";
    hash = "sha256-3vuwCOZHgmy0tV9tcgpIRjLxXa4EwNuWIbt9UkRUcDE=";
    downloadToTemp = true;
    recursiveHash = true;
    nativeBuildInputs = [ tarfs_make patchutils ];
    postFetch = ''
      # create a diff where files under fs/tarfs are added to the kernel build
      # "a" is the kernel source tree without tarfs
      # "b" is the kernel source tree with tarfs
      mkdir -p /build/a
      install -D $downloadedFile /build/b/fs/tarfs/tarfs.c
      cp -rT ${tarfs_make} /build/b
      cd /build && diff -Naur a b > /build/tarfs.patch || true
      # remove timestamps
      filterdiff --remove-timestamps /build/tarfs.patch > $out
    '';
  };
  kernel = linuxManualConfig {
    src = fetchurl {
      url = "https://cblmarinerstorage.blob.core.windows.net/sources/core/kernel-uvm-${modDirVersion}.tar.gz";
      hash = "sha256-8EU8NmU4eiqHdDeCNH28y2wKLaHx6fNcBKzWupqf2Sw=";
    };
    kernelPatches = [
      # this patches the existing Makefile and Kconfig to know about CONFIG_TARFS_FS and fs/tarfs
      { name = "build_tarfs"; patch = ./0001-kernel-uvm-6-1-build-tarfs.patch; extraConfig = { TARFS_FS = "y"; }; }
      # this adds fs/tarfs
      { name = "tarfs"; patch = tarfs_patch; }
    ];
    configfile = fetchurl {
      url = "https://raw.githubusercontent.com/microsoft/azurelinux/59ce246f224f282b3e199d9a2dacaa8011b75a06/SPECS/kernel-uvm/config";
      hash = "sha256-h13fkpQSaYnRCurkqw+zHA5BUtPxXApv6NspVAV2vXw=";
    };
    version = kver;
    modDirVersion = modDirVersion;
    # Allow reading the kernel config
    # this is required to allow nix
    # evaluation to depend on cfg
    # and correctly build everything.
    # Without this, the kernel build
    # has no support for modules.
    allowImportFromDerivation = true;
  };
  # tarfs = stdenv.mkDerivation rec {
  #   name = "tarfs-${kataVersion}-${kernel.version}";
  #   version = kataVersion;

  #   src = fetchFromGitHub {
  #     owner = "microsoft";
  #     repo = "kata-containers";
  #     rev = kataVersion;
  #     hash = "sha256-fyfPut2RFSsHZugq/zeW0+nA8F9qQNKmyhb5VqkV9Sw=";
  #   };

  #   sourceRoot = "${src.name}/src/tarfs";
  #   nativeBuildInputs = kernel.moduleBuildDependencies;

  #   makeFlags = [
  #     # setting KERNELRELEASE will attempt to add it to bzImage
  #     # like CONFIG_TARFS=y
  #     # If KERNELRELEASE is not set, tarfs will be built as a module
  #     # like CONFIG_TARFS=m
  #     # "KERNELRELEASE=${kernel.modDirVersion}"
  #     "KDIR=${kernel.dev}/lib/modules/${kernel.modDirVersion}/build"
  #     "INSTALL_MOD_PATH=$(out)"
  #   ];

  #   meta = with lib; {
  #     homepage = "https://github.com/microsoft/kata-containers";
  #     platforms = platforms.linux;
  #   };
  # };
in
kernel
