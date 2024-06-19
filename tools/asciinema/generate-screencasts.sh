#!/usr/bin/env bash
# Copyright 2024 Edgeless Systems GmbH
# SPDX-License-Identifier: AGPL-3.0-only

#
# This script prepares the environment for expect scripts to be recorded in,
# executes all scripts, and copies the .cast files to our doc's asset folder.
#

set -euo pipefail

for i in "$@"; do
  case $i in
    --demodir=*)
      demodir="${i#*=}"
      shift
      ;;
    --bin=*)
      contrastPath="${i#*=}"
      shift
      ;;
    --*)
      echo "Unknown option $i"
      exit 1
      ;;
    *)
      ;;
  esac
done

if [[ -z "${demodir:-}" ]]; then
  demodir=$(nix develop .#demo-latest --command pwd)
fi
if [[ -z "${contrastPath:-}" ]]; then
  contrastPath=$(nix build .#contrast-releases.latest && realpath result/bin/contrast)
fi
scriptdir=$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")

docker build -t screenrecodings docker

docker run -it \
  -v "${HOME}/.kube/config:/root/.kube/config" \
  -v "${scriptdir}/recordings:/recordings" \
  -v "${scriptdir}/scripts:/scripts" \
  -v "${demodir}:/demo" \
  -v "${contrastPath}:/usr/local/bin/contrast" \
  screenrecodings /scripts/flow.expect

kubectl delete -f "${demodir}/deployment/"
kubectl delete -f "${demodir}/coordinator.yml"
kubectl delete -f "${demodir}/runtime.yml"
rm result
