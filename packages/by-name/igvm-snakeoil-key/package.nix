# Copyright 2024 Edgeless Systems GmbH
# SPDX-License-Identifier: AGPL-3.0-only

{ lib
, python3
, writeTextFile
}:

python3.pkgs.buildPythonApplication rec {
  pname = "igvm-snakeoil-key";
  version = "1.5.0";
  pyproject = true;

  src = ./.;

  propagatedBuildInputs = with python3.pkgs; [
    ecdsa
    setuptools
  ];

  meta = {
    description = "Snakeoil signing key for IGVM ID block";
    mainProgram = "gen_snakeoil_pem";
    platforms = lib.platforms.all;
  };

  passthru.snakeoilPem = writeTextFile {
    name = "snakeoil.pem";
    text = builtins.readFile ./snakeoil.pem;
  };
}
