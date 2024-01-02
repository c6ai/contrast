# RFC 001: Network Sidecar Proxy

A network sidecar proxy should automatically encrypt all network traffic
between confidential applications.

## The Problem

One of the key goals of Confidential Containers is that customers can lift and
shift their existing applications to the cloud in a confidential context.

To have the full deployment secured, the communication between the pods must
also be secured and rooted in the hardware root of trust. Because of that,
the coordinator already hands out certificates which can be used for secure
(m)TLS communication.

As of today, users must change their applications source code to load
the certificates and also enable client verification inside TLS.

While users must adjust their Kubernetes deployments, they shouldn't be required
to change their application's source code as this would involve not only the
devops side but also the developer side.

## Requirements

Network traffic should be encrypted using mTLS without the need to change the
application's source code.

## Solution

Add a network sidecar proxy that establishes a service mesh between workloads
(similar to e.g. Istio).

To achieve this, we deploy a sidecar container[1] alongside the workload.

[1] <https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/>

<!--

## Alternatives Considered

### Exposing more of KubeadmConfig

We could allow users to supply their own patches to `KubeadmConfig` for finer
control over the installation. We don't want to do this because:

1. It does not solve the problem of image verification - we'd still need to
   derive image hashes from somewhere.
2. It's easy to accidentally leave charted territory when applying config
   overrides, and responsibilities are unclear in that case: should users be
   allowed to configure network, etcd, etc.?
3. The way Kubernetes exposes the configuration is an organically grown mess:
   registries are now in multiple structs, path names are hard-coded to some
   extent and versions come from somewhere else entirely (cf.
   kubernetes/kubernetes#102502).

### Ship the container images with the OS

We could bundle all control plane images in our OS image and configure kubeadm
to never pull images. This would make Constellation independent of external
image resources at the expense of flexibility: overriding the control plane
images in development setups would be harder, and we would not be able to
support user-provided images anymore.

-->
