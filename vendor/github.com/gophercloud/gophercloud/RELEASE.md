# Gophercloud release

## Contributions

### The semver label

Gophercloud follows [semver](https://semver.org/).

Each Pull request must have a label indicating its impact on the API:
* `semver:patch` for changes that don't impact the API
* `semver:minor` for changes that impact the API in a backwards-compatible fashion
* `semver:major` for changes that introduce a breaking change in the API

Automation prevents merges if the label is not present.

### Metadata

The release notes for a given release are generated based on the PR title: make
sure that the PR title is descriptive.

## Release of a new version

Requirements:
* [`gh`](https://github.com/cli/cli)
* [`jq`](https://stedolan.github.io/jq/)

### Step 1: Collect all PRs since the last release

Supposing that the base release is `v1.2.0`:

```
for commit_sha in $(git log --pretty=format:"%h" v1.2.0..HEAD); do
        gh pr list --search "$commit_sha" --state merged --json number,title,labels,url
done | jq '.[]' | jq --slurp 'unique_by(.number)' > prs.json
```

This JSON file will be useful later.

### Step 2: Determine the version

In order to determine the version of the next release, we first check that no incompatible change is detected in the code that has been merged since the last release. This step can be automated with the `gorelease` tool:

```shell
gorelease | grep -B2 -A0 '^## incompatible changes'
```

If the tool detects incompatible changes outside a `testing` package, then the bump is major.

Next, we check all PRs merged since the last release using the file `prs.json` that we generated above.

* Find PRs labeled with `semver:major`: `jq 'map(select(contains({labels: [{name: "semver:major"}]}) ))' prs.json`
* Find PRs labeled with `semver:minor`: `jq 'map(select(contains({labels: [{name: "semver:minor"}]}) ))' prs.json`

The highest semver descriptor determines the release bump.

### Step 3: Release notes and version string

Once all PRs have a sensible title, generate the release notes:

```shell
jq -r '.[] | "* [GH-\(.number)](\(.url)) \(.title)"' prs.json
```

Add that to the top of `CHANGELOG.md`. Also add any information that could be useful to consumers willing to upgrade.

**Set the new version string in the `DefaultUserAgent` constant in `provider_client.go`.**

Create a PR with these two changes. The new PR should be labeled with the semver label corresponding to the type of bump.

### Step 3: Git tag and Github release

The Go mod system relies on Git tags. In order to simulate a review mechanism, we rely on Github to create the tag through the Release mechanism.

* [Prepare a new release](https://github.com/gophercloud/gophercloud/releases/new)
* Let Github generate the  release notes by clicking on Generate release notes
* Click on **Save draft**
* Ask another Gophercloud maintainer to review and publish the release

_Note: never change a release or force-push a tag. Tags are almost immediately picked up by the Go proxy and changing the commit it points to will be detected as tampering._
