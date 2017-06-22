# kubernetes-deployment

Deploys our manifests from GitHub to Kubernetes.

## Usage

There is a complete built-in help, which is probably more up to date, than this README:

```
kubernetes-deployment help
```

### Configuration

There are different sources of retrieving parameter. They are loaded in the following order. Each item overwrites the previous one:

1. Hardcoded default values
2. Values in `~/.rebuy/kubernetes-deployment/default.[yaml|toml|json|hcl]` (eg `kubeconfig`)
2. Values in `~/.rebuy/kubernetes-deployment/$context.[yaml|toml|json|hcl]` (eg `kubeconfig`)
3. Values in `./config.[yaml|toml|json|hcl]` (eg `kubeconfig`)
4. Environment variables (eg `KUBECONFIG`)
5. Command line flags (eg `--kubeconfig`)

The key names in the configuration files always equals the long command line flags. You can generate a complete config with `kubectl dump-config`.

**Example**

```yaml
# ~/.rebuy/kubernetes-deployment/default.yaml
filename: github.com/rebuy-de/cloud-infrastructure/deployments.yaml
github-token: aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```
### Templates

All Kubernetes manifests will be rendered with the [Golang template engine](https://golang.org/pkg/text/template/).

#### Functions

The following functions are provided:

* `ToLower` - converts the string to lowercase chars
* `ToUpper` - converts the string to uppercase chars
` `Identifier` - converts the string to a valid Kubernetes identifier (eg for the `meta.name` field)

#### Variables

`kubernetes-deployment` uses variables which are inherited in a specific order. Each item overwrites the previous one:

1. generated values
2. hardcoded default values
3. default variables from project config, ie `default.variables`.
4. variables from context, ie `contexts.$context.variables`.
5. variables from service, ie `services[i].variables`.

These are the generated values:

* `gitBranchName` - The branch name from GitHub, eg `master`.
* `gitCommitID` - The full git commit hash, eg `afad13cf1941af4ad3101bdf30f087f7dfe27c99`. Useful for image tags.



