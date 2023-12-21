{ lib
, fetchFromGitHub
, buildGoModule
, mkYarnPackage
, fetchYarnDeps
, cacert
, buildContainer
, pushContainer
}:
let
  version = "0.0.1";
  src = fetchFromGitHub {
    owner = "3u13r";
    repo = "emojivoto";
    rev = "ef1e904025f2d56a925de44683a44dc418cf9aaa";
    hash = "sha256-vMyVdNu1UED/J2B3Vhrspl0y/E1SBgMd8tMnHDWIO1Y=";
  };
  vendorHash = "sha256-XuuZbExK7yXleLn16SKL6W9xPx+NgrC8RxMTlIoeZ5A=";
in
rec {
  emoji-svc = buildContainer (buildGoModule {
    inherit version src vendorHash;
    pname = "emojivoto-emoji-svc";
    subPackages = [ "emojivoto-emoji-svc/cmd" ];
    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" "-buildid=" ];
    proxyVendor = true;
    meta.mainProgram = "emojivoto-emoji-svc";
  });

  voting-svc = buildGoModule {
    inherit version src vendorHash;
    pname = "emojivoto-voting-svc";
    subPackages = [ "emojivoto-voting-svc/cmd" ];
    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" "-buildid=" ];
    proxyVendor = true;
    meta.mainProgram = "emojivoto-voting-svc";
  };

  webapp = mkYarnPackage {
    inherit version;
    pname = "emojivoto-web-webapp";

    src = "${src}/emojivoto-web/webapp";

    offlineCache = fetchYarnDeps {
      yarnLock = "${src}/emojivoto-web/webapp/yarn.lock";
      hash = "sha256-YV7amj11kEo/lgFNwYKpUpTP0eSOOqLQLVxhznx/drI=";
    };
    packageJSON = "${src}/emojivoto-web/webapp/package.json";

    buildPhase = ''
      export NODE_OPTIONS="--openssl-legacy-provider"
      yarn --offline webpack
    '';
    installPhase = ''
      mkdir -p $out
      cp deps/web/dist/index_bundle.js $out
    '';
    distPhase = "true";
  };

  web = buildGoModule rec {
    inherit version src vendorHash;
    pname = "emojivoto-web";
    subPackages = [ "emojivoto-web/cmd" ];
    CGO_ENABLED = 0;
    ldflags = [ "-s" "-w" "-buildid=" ];
    proxyVendor = true;
    passthru.ui = webapp;
    preBuild = ''
      mkdir -p emojivoto-web/ui/build
      cp -R ${passthru.ui}/ emojivoto-web/ui/build/
    '';
    meta.mainProgram = "emojivoto-web";
  };


}
