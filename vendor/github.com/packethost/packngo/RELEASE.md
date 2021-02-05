# Releases

This file serves to provide guidance and act as a checklist for any maintainers
to this project, now or in the future. This file should be updated with any
changes to the process. Automated processes should be described well enough that
they can be run in the absence of that automation.

* See [CHANGELOG.md](CHANGELOG.md) for notes on versioning.
* Fetch the latest origin branches:

  ```sh
  git fetch origin
  git checkout master
  git pull
  ```

* Verify that your branch matches the upstream branch:

  ```sh
  git branch --points-at=master -r  | grep origin/master >/dev/null || echo "master differs from origin/master"
  ```

* Update the `libraryVersion` constant. This is a library, so we can not assure
  that a build flag will be used in every client that provides a compile time
  value, let alone the correct one.

  ```sh
  vim packngo.go # change libraryVersion, "0.3.0" (no v)
  git commit --signoff -m 'v0.3.0 version bump' packngo.go
  ```

* Tag `master` with a semver tag that suits the level of changes
  introduced:

  ```sh
  git tag -m "v0.3.0" -a v0.3.0 master # use -s if gpg is available
  ```
* Push the tag:

  ```sh
  git push --tags origin master v0.3.0
  ```
* Create a release from the tag (include a keepthechangelog.com formatted description):

  <https://github.com/packethost/packngo/releases/edit/v0.3.0> (use the correct
  version)

Releases can be followed through the GitHub Atom feed at
<https://github.com/packethost/packngo/releases.atom>.
