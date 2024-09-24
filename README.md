# hellas

## Developing Cycle

## Requirements

* Minikube
* Skaffold

### Getting Started

Run `make up` to stand up a minikube cluster and start Skaffold

### Cleanup

Run `make clean` to tear down and clean up minikube and stop Skaffold

```yaml
server:
  certSecretName: my-cert

registry:
  github:
    tokenSecretName: my-token
    insecureSkipVerify: false
    protocol: https
  # s3:
  #   region: us-east-1
  #   bucket: foo
  #   path: bar/bix
```

TODO:

* set up cert manager integration
* set up logging
* add k8s health checks
