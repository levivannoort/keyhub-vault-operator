<p align="center">
  <img src="assets/keyhub.png" width="345" height="180">
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/topicuskeyhub/keyhub-vault-operator" alt="release">
  <img src="https://img.shields.io/github/go-mod/go-version/topicuskeyhub/keyhub-vault-operator" alt="go">
  <img src="https://img.shields.io/github/license/topicuskeyhub/keyhub-vault-operator" alt="license">
</p>

# Topicus - KeyHub Vault Operator

The KeyHub Vault Operator can be used in conjuction with a Topicus KeyHub instance, to be able to synchronize records in a KeyHub vault to Kubernetes secrets. It does so through a policy based mechanism, which are stored in a KeyHub vault. 

To define a mapping between KeyHub and Kubernetes a `KeyHubSecret` custom resource can be created. The name of the generated Kubernetes Secret is the same as the name of the KeyHubSecret. The mapping between a secret key and a vault record is based on the uuid of the vault record. Secrets will be automatically synchronized with a 10 minute interval. In case of an error the retry interval is 2 minutes.

> Documentation: https://topicuskeyhub.github.io/keyhub-vault-operator/
>
> Product site: https://www.topicus-keyhub.com/

## Installation

Although the KeyHub Vault Operator can be installed within the kubernetes cluster without the pre-requisites, it wont function if the following things aren't present:
- Accessible KeyHub instance.
- Configuration of KeyHub components (e.g. vaults, OIDC applications, policy, records).
- Kubernetes secret with the client credentials (i.e. `keyhub-vault-operator-secret`), used to access the policy vault.

The following kustomize file can be applied to the Kubernetes cluster for the installation of the controller.
```
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: keyhub-vault-operator

resources:
- ssh://github.com/topicuskeyhub/keyhub-vault-operator//config/default?ref=main
```
## Configuration

- []()
- []()

## Release
Manually run the `release` workflow (branch `main`) from Github Actions. This will create a release which in turn is the trigger for the `build-and-publish` workflow.

The `release` workflow will do the following:
- Update `images['controller'].newTag` in `config/manager/kustomization.yaml` with the full semver (prefixed with a 'v'), e.g. `v0.1.0`. The semver is based on conventional commits and the latest git tag.
- Create and push the release tag, e.g. `v0.1.0`.
- Create a GitHub release with a changelog based on the (conventional) commits since the last release.

The `build-and-publish` workflow will do the following.
- Build and publish the image from the release tag.

## Development - requirements

- Installation of a local kubernetes cluster, e.g. [minikube](https://minikube.sigs.k8s.io/docs/)
- Install [Operator SDK](https://sdk.operatorframework.io/)

### Pre-commit hook to check conventional commits
- Install [`pre-commit`](https://pre-commit.com/#install)
- Configure [`pre-commit-script`](https://github.com/compilerla/conventional-pre-commit)

```console
pre-commit install --hook-type commit-msg
```

## Development - getting started

Run the operator locally. Make sure you are connecting to your local Kubernetes cluster!
```
make run
```

Run the tests
```
make test
```