# Troubleshooting

This section contains information on how to debug your Contrast deployment.

## Logging

Collecting logs can be a good first step to identify problems in your
deployment. Both the CLI and the Contrast Coordinator as well as the Initializer
can be configured to emit additional logs.

### CLI

The CLI logs can be configured with the `--log-level` command-line flag, which
can be set to either `debug`, `info`, `warn` or `error`. The default is `info`.
Setting this to `debug` can get more fine-grained information as to where the
problem lies.

### Coordinator and Initializer

The logs from the Coordinator and the Initializer can be configured via the
environment variables `CONTRAST_LOG_LEVEL`, `CONTRAST_LOG_FORMAT` and
`CONTRAST_LOG_SUBSYSTEMS`.
- `CONTRAST_LOG_LEVEL` can be set to one of either `debug`, `info`, `warn`, or
  `error`, similar to the CLI (defaults to `info`).
- `CONTRAST_LOG_FORMAT` can be set to `text` or `json`, determining the output
  format (defaults to `text`).
- `CONTRAST_LOG_SUBSYSTEMS` is a comma-seperated list of subsystems that should
  be enabled for logging, which are disabled by default. Subsystems include:
  `snp-issuer`, `kds-getter`, and `snp-validator`. To enable all subsystems, use
  `*` as the value for this environment variable.

Warnings and error messages from subsystems get printed regardless of whether
the subsystem is listed in the `CONTRAST_LOG_SUBSYSTEMS` environment variable.

#### Coordinator debug logging

To configure debug logging with all subsystems for your Coordinator, add the
following variables to your container definition.

```yaml
spec: # v1.PodSpec
  containers:
    image: "ghcr.io/edgelesssys/contrast/coordinator:latest"
    name: coordinator
    env:
    - name: CONTRAST_LOG_LEVEL
      value: debug
    - name: CONTRAST_LOG_SUBSYSTEMS
      value: "*"
    # ...
```

To access the logs generated by the Coordinator, you can use `kubectl` with the
following command.

```sh
kubectl logs <coordinator-pod-name>
```

#### Pod fails on startup

If the Coordinator or a workload pod fails to even start, it can be helpful to
look at the events of the pod during the startup process using the `describe`
command.

```sh
kubectl describe pod <pod-name>
```

Example output:

```
Events:
  Type     Reason   Age    From     Message
  ----     ------   ----   ----     -------
  ...
  Warning  Failed   20s    kubelet  Error: failed to create containerd task: failed to create shim task: "CreateContainerRequest is blocked by policy: ...
```

In this example, the container creation was blocked by a policy. This suggests
that a policy hasn't been updated to accommodate recent changes to the
configuration. Make sure to run `contrast generate` when altering your
deployment.
