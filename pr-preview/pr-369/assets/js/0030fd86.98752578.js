"use strict";(self.webpackChunkcontrast_docs=self.webpackChunkcontrast_docs||[]).push([[8242],{1289:e=>{e.exports=JSON.parse('{"pluginId":"default","version":"0.5.0","label":"0.5.0","banner":null,"badge":true,"noIndex":false,"className":"docs-version-0.5.0","isLast":true,"docsSidebars":{"docs":[{"type":"category","label":"What is Contrast?","collapsed":false,"items":[{"type":"link","label":"Confidential Containers","href":"/contrast/pr-preview/pr-369/basics/confidential-containers","docId":"basics/confidential-containers","unlisted":false},{"type":"link","label":"Security benefits","href":"/contrast/pr-preview/pr-369/basics/security-benefits","docId":"basics/security-benefits","unlisted":false},{"type":"link","label":"Features","href":"/contrast/pr-preview/pr-369/basics/features","docId":"basics/features","unlisted":false}],"collapsible":true,"href":"/contrast/pr-preview/pr-369/"},{"type":"category","label":"Getting started","collapsed":false,"items":[{"type":"link","label":"Install","href":"/contrast/pr-preview/pr-369/getting-started/install","docId":"getting-started/install","unlisted":false},{"type":"link","label":"Cluster setup","href":"/contrast/pr-preview/pr-369/getting-started/cluster-setup","docId":"getting-started/cluster-setup","unlisted":false},{"type":"link","label":"First steps","href":"/contrast/pr-preview/pr-369/getting-started/first-steps","docId":"getting-started/first-steps","unlisted":false}],"collapsible":true,"href":"/contrast/pr-preview/pr-369/getting-started/"},{"type":"category","label":"Examples","items":[{"type":"link","label":"Confidential emoji voting","href":"/contrast/pr-preview/pr-369/examples/emojivoto","docId":"examples/emojivoto","unlisted":false}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/examples/"},{"type":"link","label":"Workload deployment","href":"/contrast/pr-preview/pr-369/deployment","docId":"deployment","unlisted":false},{"type":"category","label":"Architecture","items":[{"type":"category","label":"Components","items":[{"type":"link","label":"Coordinator","href":"/contrast/pr-preview/pr-369/architecture/components/coordinator","docId":"architecture/components/coordinator","unlisted":false},{"type":"link","label":"Init container","href":"/contrast/pr-preview/pr-369/architecture/components/init-container","docId":"architecture/components/init-container","unlisted":false},{"type":"link","label":"CLI","href":"/contrast/pr-preview/pr-369/architecture/components/cli","docId":"architecture/components/cli","unlisted":false}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/category/components"},{"type":"link","label":"Confidential Containers","href":"/contrast/pr-preview/pr-369/architecture/confidential-containers","docId":"architecture/confidential-containers","unlisted":false},{"type":"category","label":"Attestation","items":[{"type":"link","label":"Hardware","href":"/contrast/pr-preview/pr-369/architecture/attestation/hardware","docId":"architecture/attestation/hardware","unlisted":false},{"type":"link","label":"Pod VM","href":"/contrast/pr-preview/pr-369/architecture/attestation/pod-vm","docId":"architecture/attestation/pod-vm","unlisted":false},{"type":"link","label":"Runtime policies","href":"/contrast/pr-preview/pr-369/architecture/attestation/runtime-policies","docId":"architecture/attestation/runtime-policies","unlisted":false},{"type":"link","label":"Manifest","href":"/contrast/pr-preview/pr-369/architecture/attestation/manifest","docId":"architecture/attestation/manifest","unlisted":false},{"type":"link","label":"Coordinator","href":"/contrast/pr-preview/pr-369/architecture/attestation/coordinator","docId":"architecture/attestation/coordinator","unlisted":false}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/category/attestation"},{"type":"category","label":"Certificates and Identities","items":[{"type":"link","label":"PKI","href":"/contrast/pr-preview/pr-369/architecture/certificates-and-identities/pki","docId":"architecture/certificates-and-identities/pki","unlisted":false}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/category/certificates-and-identities"},{"type":"category","label":"Network Encryption","items":[{"type":"link","label":"Sidecar","href":"/contrast/pr-preview/pr-369/architecture/network-encryption/sidecar","docId":"architecture/network-encryption/sidecar","unlisted":false},{"type":"link","label":"Protocols and Keys","href":"/contrast/pr-preview/pr-369/architecture/network-encryption/protocols-and-keys","docId":"architecture/network-encryption/protocols-and-keys","unlisted":false}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/category/network-encryption"}],"collapsed":true,"collapsible":true,"href":"/contrast/pr-preview/pr-369/architecture/"}]},"docs":{"architecture/attestation/coordinator":{"id":"architecture/attestation/coordinator","title":"coordinator","description":"","sidebar":"docs"},"architecture/attestation/hardware":{"id":"architecture/attestation/hardware","title":"hardware","description":"","sidebar":"docs"},"architecture/attestation/manifest":{"id":"architecture/attestation/manifest","title":"manifest","description":"","sidebar":"docs"},"architecture/attestation/pod-vm":{"id":"architecture/attestation/pod-vm","title":"pod-vm","description":"","sidebar":"docs"},"architecture/attestation/runtime-policies":{"id":"architecture/attestation/runtime-policies","title":"runtime-policies","description":"","sidebar":"docs"},"architecture/certificates-and-identities/pki":{"id":"architecture/certificates-and-identities/pki","title":"pki","description":"","sidebar":"docs"},"architecture/components/cli":{"id":"architecture/components/cli","title":"cli","description":"","sidebar":"docs"},"architecture/components/coordinator":{"id":"architecture/components/coordinator","title":"coordinator","description":"","sidebar":"docs"},"architecture/components/init-container":{"id":"architecture/components/init-container","title":"init-container","description":"","sidebar":"docs"},"architecture/confidential-containers":{"id":"architecture/confidential-containers","title":"confidential-containers","description":"","sidebar":"docs"},"architecture/index":{"id":"architecture/index","title":"Architecture","description":"","sidebar":"docs"},"architecture/network-encryption/protocols-and-keys":{"id":"architecture/network-encryption/protocols-and-keys","title":"protocols-and-keys","description":"","sidebar":"docs"},"architecture/network-encryption/sidecar":{"id":"architecture/network-encryption/sidecar","title":"sidecar","description":"","sidebar":"docs"},"basics/confidential-containers":{"id":"basics/confidential-containers","title":"Confidential Containers","description":"Contrast uses some building blocks from Confidential Containers (CoCo), a CNCF Sandbox project that aims to standardize confidential computing at the pod level.","sidebar":"docs"},"basics/features":{"id":"basics/features","title":"Product Features","description":"Contrast simplifies the deployment and management of Confidential Containers, offering optimal data security for your workloads while integrating seamlessly with your existing Kubernetes environment.","sidebar":"docs"},"basics/security-benefits":{"id":"basics/security-benefits","title":"security-benefits","description":"","sidebar":"docs"},"deployment":{"id":"deployment","title":"Workload deployment","description":"The following instructions will guide you through the process of making an existing Kubernetes deployment","sidebar":"docs"},"examples/emojivoto":{"id":"examples/emojivoto","title":"Confidential emoji voting","description":"screenshot of the emojivoto UI","sidebar":"docs"},"examples/index":{"id":"examples/index","title":"Examples","description":"","sidebar":"docs"},"getting-started/cluster-setup":{"id":"getting-started/cluster-setup","title":"Create a cluster","description":"Prerequisites","sidebar":"docs"},"getting-started/first-steps":{"id":"getting-started/first-steps","title":"first-steps","description":"","sidebar":"docs"},"getting-started/index":{"id":"getting-started/index","title":"Getting started","description":"","sidebar":"docs"},"getting-started/install":{"id":"getting-started/install","title":"Installation","description":"Download the bundle from the URL you received:","sidebar":"docs"},"intro":{"id":"intro","title":"Contrast","description":"Welcome to the documentation of Contrast! Contrast runs confidential container deployments on Kubernetes at scale.","sidebar":"docs"}}}')}}]);