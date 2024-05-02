{ lib
, stdenv
, dnf-plugins-core
, writeText
, writeTextDir
, writeShellApplication
, dnf4
, jq
, wget
, symlinkJoin
, nix
}:
let
  dnfConf = writeText "dnf.conf" ''
    [main]
    gpgcheck=True
    installonly_limit=3
    clean_requirements_on_remove=True
    best=False
    skip_if_unavailable=True
    tsflags=nodocs
    pluginpath=${dnf-plugins-core}/lib/python3.11/site-packages/dnf-plugins
  '';
  reposdir = writeTextDir "yum.repos.d/cbl-mariner-2.repo" ''
    [cbl-mariner-2.0-prod-base-x86_64-yum]
    name=cbl-mariner-2.0-prod-base-x86_64-yum
    baseurl=https://packages.microsoft.com/yumrepos/cbl-mariner-2.0-prod-base-x86_64/
    repo_gpgcheck=0
    gpgcheck=0
    enabled=1
    gpgkey=https://packages.microsoft.com/yumrepos/cbl-mariner-2.0-prod-base-x86_64/repodata/repomd.xml.key
  '';
  packages = writeText "packages" ''
    kata-packages-uvm
    kata-packages-uvm-coco
    systemd
    libseccomp
    opa
  '';
  update_lockfile = writeShellApplication {
    name = "update_lockfile";
    runtimeInputs = [ dnf4 jq wget nix ];
    text = builtins.readFile ./update_lockfile.sh;
  };
in
symlinkJoin {
  name = "rpm-pin-vendor";
  paths = [ ];
  postBuild = ''
    mkdir -p $out/bin
    cp ${lib.getExe update_lockfile} $out/bin/update_lockfile
    substituteInPlace $out/bin/update_lockfile \
      --replace-fail "@DNFCONFIG@" ${dnfConf} \
      --replace-fail "@REPOSDIR@" ${reposdir}/yum.repos.d \
      --replace-fail "@PACKAGESET@" ${packages}
  '';
  meta.mainProgram = "update_lockfile";
}
