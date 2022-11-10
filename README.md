# dc2deploy

```Convert Openshift DeploymentConfig to Kuberentes Deployment

Usage:
  dc2deploy [flags]

Flags:
      --dry-run             Only print the new object that would be sent
  -f, --filename string     File containing DeploymentConfig manifest (default "-")
  -h, --help                help for dc2deploy
      --ignore-warnings     Ignore Warnings about missing Deployment Features
      --kubeconfig string   Path to Kubeconfig
  -n, --namespace string    Namespace of DeploymentConfig
      --outfile string      Output filename. Defaults to STDOUT (default "-")
  -o, --output string       Output in JSON (default "yaml")
  -v, --v Level             number for the log level verbosity
```
