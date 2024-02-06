{
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , treefmt-nix
    , ...
    }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
      inherit (pkgs) lib;

      version = "0.2.0-pre";

      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;

      packages = import ./packages { inherit pkgs version; };
      packages_x86_64-linux = import ./packages { inherit version; pkgs = pkgs.pkgsCross.x86_64-linux; };
      packages_aarch64-darwin = import ./packages { inherit version; pkgs = pkgs.pkgsCross.aarch64-darwin; };
      packages_aarch64-multiplatform = import ./packages { inherit version; pkgs = pkgs.pkgsCross.aarch64-multiplatform; };
      packages_x86_64-darwin = import ./packages { inherit version; pkgs = pkgs.pkgsCross.x86_64-darwin; };
    in
    {
      packages = packages //
        {
          pkgsCross = lib.recurseIntoAttrs {
            aarch64-darwin = packages_aarch64-darwin;
            aarch64-multiplatform = packages_aarch64-multiplatform;
            x86_64-darwin = packages_x86_64-darwin;
          };
        };

      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          golangci-lint
          just
        ];
        shellHook = ''alias make=just'';
      };

      formatter = treefmtEval.config.build.wrapper;

      checks = {
        formatting = treefmtEval.config.build.check self;
      };

      # legacyPackages = pkgs;
    });

  nixConfig = {
    extra-substituters = [
      "https://edgelesssys.cachix.org"
    ];
    extra-trusted-public-keys = [
      "edgelesssys.cachix.org-1:erQG/S1DxpvJ4zuEFvjWLx/4vujoKxAJke6lK2tWeB0="
    ];
  };
}
