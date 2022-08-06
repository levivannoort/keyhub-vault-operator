## Development
- Local instance of a Kubernetes cluster (e.g. [minikube](https://minikube.sigs.k8s.io/docs/))
- Install [operator SDK](https://sdk.operatorframework.io/)

## Getting started
Run the operator locally. Make sure you are connecting to your local Kubernetes cluster!
```
make run
```

Run the tests
```
make test
```

### Git pre-commit hook to check Conventional Commits
- Install [`pre-commit`](https://pre-commit.com/#install)
- Install `pre-commit` script ([more info](https://github.com/compilerla/conventional-pre-commit)):

  ```console
  pre-commit install --hook-type commit-msg
  ```