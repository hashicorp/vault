# How to create a `govmomi` Release on Github

> **Note**
>
> The steps outlined in this document can only be performed by maintainers or
> administrators of this project.

The release automation is based on Github
[Actions](https://github.com/features/actions) and has been improved over time
to simplify the experience for creating `govmomi` releases.

The Github Actions release [workflow](.github/workflows/govmomi-release.yaml)
uses [`goreleaser`](http://goreleaser.com/) and automatically creates/pushes:

- Release artifacts for `govc` and `vcsim` to the
  [release](https://github.com/vmware/govmomi/releases) page, including
  `LICENSE.txt`, `README` and `CHANGELOG`
- Docker images for `vmware/govc` and `vmware/vcsim` to Docker Hub
- Source code

Releases are not tagged on the `main` branch, but a dedicated release branch, for example `release-0.35`.

### Verify `main` branch is up to date with the remote

```console
git checkout main
git fetch -avp
git diff main origin/main

# if your local and remote branches diverge run
git pull origin/main
```

> **Warning**
>
> These steps assume `origin` to point to the remote
> `https://github.com/vmware/govmomi`, respectively
> `git@github.com:vmware/govmomi`.

### Create a release branch

For new releases, create a release branch from the most recent commit in
`main`, e.g. `release-0.35`.

```console
export RELEASE_BRANCH=release-0.35
git checkout -b ${RELEASE_BRANCH}
```

For maintenance/patch releases on **existing** release branches, simply checkout the existing
release branch and add commits to the existing release branch.

### Verify `make docs` and `CONTRIBUTORS` are up to date

> **Warning**
>
> Run the following commands and commit any changes to the release branch before
> proceeding with the release.

```console
make doc
./scripts/contributors.sh
if [ -z "$(git status --porcelain)" ]; then
  echo "working directory clean: proceed with release"
else
  echo "working directory dirty: please commit changes"
fi

# perform git add && git commit ... in case there were changes
```

### Push the release branch

> **Warning**
>
> Do not create a tag as this will be done by the release automation.

The final step is pushing the new/updated release branch.

```console
git push origin ${RELEASE_BRANCH}
```

### Create a release in the Github UI

Open the `govmomi` Github [repository](https://github.com/vmware/govmomi) and
navigate to `Actions -> Workflows -> Release`.

Click `Run Workflow` which opens a dropdown list.

Select the new/updated branch, e.g. `release-0.35`, i.e. **not** the `main`
branch.

Specify a semantic `tag` to associate with the release, e.g. `v0.35.0`.

> **Warning**
>
> This tag **must not** exist or the release will fail during the validation
> phase.

By default, a dry-run is performed to rule out most (but not all) errors during
a release. If you do not want to perform a dry-run, e.g. to finally create a
release, deselect the `Verify release workflow ...` checkbox.

Click `Run Workflow` to kick off the workflow.

After successful completion and if the newly created `tag` is the **latest**
(semantic version sorted) tag in the repository, a PR is automatically opened
against the `main` branch to update the `CHANGELOG`.
