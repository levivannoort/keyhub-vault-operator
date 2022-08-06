<p align="center">
  <img src="assets/keyhub.png" width="345" height="180">
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/topicuskeyhub/keyhub-vault-operator" alt="release">
  <img src="https://img.shields.io/github/go-mod/go-version/topicuskeyhub/keyhub-vault-operator" alt="go">
  <img src="https://img.shields.io/github/license/topicuskeyhub/keyhub-vault-operator" alt="license">
</p>

# Topicus KeyHub Vault Operator

The KeyHub Vault Operator can be used in conjuction with a Topicus KeyHub instance, to be able to synchronize records in a KeyHub vault to Kubernetes secrets with the use of the `KeyhubSecret` custom resource. It does so through a policy base mechanism 

The documentation can be found [here](https://topicuskeyhub.github.io/keyhub-vault-operator/). Additional information about development/contribution and releasing can be found under the 'docs'-directory.

## Installation

Although the KeyHub Vault Operator can be installed within the Kubernetes cluster without the pre-requisites, it wont function if the following things aren't present:
- Accessible KeyHub instance.
- Configuration of KeyHub components (e.g. vaults, OIDC applications, policy record, records).
- Kubernetes secret with the client credentials (i.e. `keyhub-vault-controller-secret`).

The following kustomize file can be applied to the Kubernetes cluster for the installation of the controller.
```
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: keyhub-vault-operator

resources:
- ssh://github.com/topicuskeyhub/keyhub-vault-operator//config/default?ref=main
```
