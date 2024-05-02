{ lib
, stdenv
, stdenvNoCC
, distro ? "cbl-mariner"
, bubblewrap
, fakeroot
, fetchFromGitHub
, fetchurl
, kata-agent
, yq-go
, tdnf
, curl
, util-linux
, writeText
, writeTextDir
, createrepo_c
, qemu-utils
, writeShellApplication
, e2fsprogs
, parted
, cryptsetup
}:
# https://github.com/microsoft/azurelinux/blob/59ce246f224f282b3e199d9a2dacaa8011b75a06/SPECS/kata-containers-cc/mariner-coco-build-uvm.sh#L18
let
  kata-version = "3.2.0.azl1";
  source = fetchFromGitHub {
    owner = "microsoft";
    repo = "kata-containers";
    rev = kata-version;
    hash = "sha256-fyfPut2RFSsHZugq/zeW0+nA8F9qQNKmyhb5VqkV9Sw=";
  };
  packageIndex = builtins.fromJSON (builtins.readFile ./package-index.json);
  rpmSources = lib.forEach packageIndex
    (p: lib.concatStringsSep "#" [ (fetchurl p) (builtins.baseNameOf p.url) ]);

  mirror = stdenvNoCC.mkDerivation {
    name = "mirror";
    dontUnpack = true;
    nativeBuildInputs = [ createrepo_c ];
    buildPhase = ''
      runHook preBuild

      mkdir -p $out/packages
      for source in ${builtins.concatStringsSep " " rpmSources}; do
        path=$(echo $source | cut -d'#' -f1)
        filename=$(echo $source | cut -d'#' -f2)
        ln -s "$path" "$out/packages/$filename"
      done

      createrepo_c --revision 0 --set-timestamp-to-revision --basedir packages $out

      runHook postBuild
    '';
  };

  tdnfConf = writeText "tdnf.conf" ''
    [main]
    gpgcheck=1
    installonly_limit=3
    clean_requirements_on_remove=0
    repodir=/etc/yum.repos.d
    cachedir=/build/var/cache/tdnf
  '';
  vendor-reposdir = writeTextDir "yum.repos.d/cbl-mariner-2-vendor.repo" ''
    [cbl-mariner-2.0-prod-base-x86_64-yum]
    name=cbl-mariner-2.0-prod-base-x86_64-yum
    baseurl=file://${mirror}
    repo_gpgcheck=0
    gpgcheck=0
    enabled=1
  '';
  rootfs = stdenv.mkDerivation rec {
    pname = "kata-rootfs";
    version = kata-version;

    env = {
      AGENT_SOURCE_BIN = "${lib.getExe kata-agent}";
      # TODO: understand why build fails with AGENT_POLICY enabled
      # AGENT_POLICY = "yes";
      CONF_GUEST = "yes";
      # TODO: Add support for custom policy file
      # AGENT_POLICY_FILE=allow-set-policy.rego
      RUST_VERSION = "not-needed";
    };

    nativeBuildInputs = [
      yq-go
      curl
      fakeroot
      bubblewrap
      util-linux
      tdnf
    ];

    src = source;

    sourceRoot = "${src.name}/tools/osbuilder/rootfs-builder";

    buildPhase = ''
      runHook preBuild

      mkdir -p /build/var/run
      mkdir -p /build/var/tdnf
      mkdir -p /build/var/lib/tdnf
      mkdir -p /build/var/cache/tdnf
      mkdir -p /build/root
      unshare --map-root-user bwrap \
        --bind /nix /nix \
        --bind ${tdnfConf} /etc/tdnf/tdnf.conf \
        --bind ${vendor-reposdir}/yum.repos.d /etc/yum.repos.d \
        --bind /build /build \
        --bind /build/var /var \
        --dev-bind /dev/null /dev/null \
        fakeroot bash -c "bash $(pwd)/rootfs.sh -r /build/root ${distro} && \
          tar --sort=name --mtime='UTC 1970-01-01' -C /build/root -c . -f $out"

      runHook postBuild
    '';
  };
  buildimage = writeShellApplication {
    name = "buildimage";
    runtimeInputs = [ e2fsprogs cryptsetup ];
    text = builtins.readFile ./buildimage.sh;
  };
in
stdenvNoCC.mkDerivation rec {
  pname = "kata-image";
  version = kata-version;

  nativeBuildInputs = [
    qemu-utils
    fakeroot
    bubblewrap
    util-linux
    parted
    buildimage
  ];

  dontUnpack = true;

  buildPhase = ''
    runHook preBuild

    mkdir -p /build/rootfs
    unshare --map-root-user bwrap \
      --bind /nix /nix \
      --bind /build /build \
      --bind /build/rootfs /rootfs \
      --dev-bind /dev/null /dev/null \
      --dev-bind /dev/random /dev/random \
      --dev-bind /dev/urandom /dev/urandom \
      fakeroot bash -c "tar -pxf ${rootfs} -C /rootfs  && bash ${lib.getExe buildimage} /rootfs $out"

    runHook postBuild
  '';

  passthru.rootfs = rootfs;
}
