# kubernetes-deployment
Deploys our Kubernetes manifests to a fresh cluster

## Usage

You can discover all flags with the tool itself:

```bash
kubernetes-deployment help
```

### Examples

```bash
kubernetes-deployment deploy -g fetch -g render -n message-broker -s production
```

```bash
kubernetes-deployment deploy -g fetch -g render -n message-broker -b cloud-42
```

```bash
kubernetes-deployment bulk -g all
```
