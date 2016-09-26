# kubernetes-deployment
Deploys our Kubernetes manifests to a fresh cluster

## Usage

You can discover all flags with the tool itself:

```bash
kubernetes-deployment -h
```

Default settings from service.yaml 
```yaml
    settings:
      kubeconfig: ~/.kube/kubeconfig.yml
      output: ./output
      sleep: 1s
      retry-sleep: 250ms
      retry-count: 3
      template-values:
        - key: clusterDomain
          value: main.cloud.rebuy.loc
    services:
    # k8s manifests locations
```

All of the settings values can be overwritten by local service.yaml file
eg ~/sandbox.yml could look like:

```yaml
    settings:
      kubeconfig: ~/.kube/sandbox.yml
      output: ./output-sandbox
      sleep: 5s      
      template-values:
        - key: clusterDomain
          value: sandbox.cloud.rebuy.loc
```

A typical usage would be:

```bash
kubernetes-deployment \
	-config ~/sandbox.yml \
    -ignore-deploy-failures \
# The `-config` flag actually default to these values.
```

This run generates reproducible file, so you are able to apply the exact same
Kubernetes manifest in the same order:

```bash
kubernetes-deployment \
    -config ./output/config.yml \
    -ignore-deploy-failures \
    -skip-fetch \
    -skip-shuffle
```

Note that `-config` changend and `-skip-shuffle` was added to apply the
manifests in the same order. Also `-skip-fetch` was set, so the tool doesn't
fetch new manifests from git.
