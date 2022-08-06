## Release
Manually run the 'release'-workflow (branch `main`) from Github Actions. The 'release'-workflow will do the following:
- Update `images['controller'].newTag` in `config/manager/kustomization.yaml` with the full semver (prefixed with a 'v'), e.g. `v0.1.0`. The semver is based on conventional commits and the latest git tag.
- Create and push the release tag, e.g. `v0.1.0`.
- Create a release with a changelog based on the (conventional) commits since the last release.
- Github actions will build and publish the image from the release tag.