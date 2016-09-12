# kubernetes-deployment
Deploys our Kubernetes manifests to a fresh cluster

## Usage

You can discover all flags with the tool itself:

```bash
kubernetes-deployment -h
```

A typical usage would be:

```bash
kubernetes-deployment \
	-config ./config/services.yaml \
    -kubeconfig ~/dev/rebuy/terraform-templates/outputs/kubeconfig.yml \
    -ignore-deploy-failures \
    -output ./output
# The `-config` and `-output` flags actually default to these values.
```

This run generates reproducible file, so you are able to apply the exact same
Kubernetes manifest in the same order:

```bash
kubernetes-deployment \
    -config ./output/services.yaml \
    -kubeconfig ~/dev/rebuy/terraform-templates/outputs/kubeconfig.yml \
    -ignore-deploy-failures \
    -output ./output \
    -skip-fetch \
    -skip-shuffle
```

Note that `-config` changend and `-skip-shuffle` was added to apply the
manifests in the same order. Also `-skip-fetch` was set, so the tool doesn't
fetch new manifests from git.
