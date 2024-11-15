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

* Tag `master` with a semver tag that suits the level of changes
  introduced:

  ```sh
  git tag -m "v0.8.0" -a v0.8.0 master # use -s if gpg is available
  ```
* Push the tag:

  ```sh
  git push --tags origin master v0.8.0
  ```
* Create a release from the tag (include a keepachangelog.com formatted description):

  <https://github.com/packethost/packngo/releases/edit/v0.8.0> (use the correct
  version)

Releases can be followed through the GitHub Atom feed at
<https://github.com/packethost/packngo/releases.atom>.
