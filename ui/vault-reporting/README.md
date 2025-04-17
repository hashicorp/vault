# Vault Reporting

This repo contains shared UI components for Vault Enterprise and Cloud reporting views.

## Publishing

At this time this packages is not published to the NPM registry. Since this repo is private and all enterprise features are built into the public open source vault repo we need a way to "publish" updates to vault. There are a couple of utilities to help with tihs.

### Github workflow

When your changes are ready and you want to consume them in vault you can trigger a github action to automatically open a PR with the updates in the vault repo.

1. Visit the [publish-reporting-to-vault](https://github.com/hashicorp/shared-secure-ui/actions/workflows/publish-reporting-to-vault.yml) workflow for `shared-secure-ui` repo.
2. Click the "Run workflow" button in the top right
3. Enter the branch in vault you want the PR to go into, by default it will be main
4. Submit the form

This should trigger the workflow and after a short amount of time you should be a new draft PR in the vault repo with the `vault-reporting` updates.

### Publishing locally

You will need to have this repo and the vault repo checked out locally.

1. Build the addon `npm run build`
2. Export an environment variable called `VAULT_UI_PATH` with the path to your vault/ui directory (if unset it will try to find it at `~/projects/vault/ui/`)
3. Run the sync script `npm run sync-to-vault --workspace @hashicorp/vault-reporting`

You should see the updated dist files in the vault-reporting directory inside of vault.
