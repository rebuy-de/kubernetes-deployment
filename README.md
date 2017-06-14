# kubernetes-deployment

Deploys our manifests from GitHub to Kubernetes.

## Usage

There is a complete built-in help, which is probably more up to date, than this README:

```
kubernetes-deployment help
```

### Configuration

There are different sources of retrieving parameter. They are loaded in the following order. Each item overwrites the previous one:

1. default values
2. values in `~/.rebuy/kubernetes-deployment/default.[yaml|toml|json|hcl]` (eg `kubeconfig`)
3. values in `./config.[yaml|toml|json|hcl]` (eg `kubeconfig`)
4. environment variables (eg `KUBECONFIG`)
5. command line flags (eg `--kubeconfig`)
