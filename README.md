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

Although the KeyHub Vault Operator can be installed within the Kubernetes cluster without the pre-requisites, it wont function if the following things aren't present:
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
