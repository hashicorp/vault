# Contributing

In addition to filing bugs, you may contribute by submitting patches to fix bugs in the library.
Contributions may be submitting to <http://review.couchbase.com>.
We use Gerrit as our code review system - and thus submitting a change requires an account there.
Note that pull requests will not be ignored but will be responded to more quickly and with more detail in Gerrit.

For something to be accepted into the codebase, it must be formatted properly and have undergone proper testing.
We use `golangci` for linting with a number of linters enabled.
To install and use the linter you can use `make devsetup` and `make lint`.
You can also run the linter and test suite together using `make check`.
Please note that we keep the linting tools out of the `go.mod` file.

## Branches and Tags

Released versions of the library are marked as annotated tags inside the repository.

* The `master` branch represents the mainline branch. The master branch typically consists of content going into the next release.

## Contributing Patches

If you wish to contribute a new feature or a bug fix to the library, try to follow the following guidelines to help
ensure your change gets merged upstream.

### Before you begin

For any code change, ensure the new code you write looks similar to the code surrounding it and that linting does not produce errors.

If your change is going to involve a substantial amount of time or effort, please attempt to discuss it with the project
developers first who will provide assistance and direction where possible.

#### For new features

Ensure the feature you are adding does not already exist, and think about how this feature may be useful for other users.
In general less intrusive changes are more likely to be accepted.

#### For fixing bugs

Ensure the bug you are fixing is actually a bug (and not a usage error), and that it has not been fixed in a more recent
version. Please read the release notes as well as the issue tracker to see a list of open and resolved issues.

### Code Review

#### Signing up on Gerrit

Everything that is merged into the library goes through a code review process.  The code review process is done via
[Gerrit](http://review.couchbase.org).

To sign up for a gerrit account, go to http://review.couchbase.org and click on the _Register_ link at the top right.
Once you've signed in you will need to agree to the CLA (Contributor License Agreement) by going you your gerrit
account page and selecting the _Agreements_ link on the left.
When you've done that, everything should flow through just fine.
Be sure that you have registered your email address at http://review.couchbase.org/#/settings/contact as many sign-up methods
won't pass emails along.
Note that your email address in your code commit and in the gerrit settings must match.

Add your public SSH key to gerrit before submitting.

#### Setting up your fork with Gerrit

Assuming you have a repository created like so:

```
$ git clone git://github.com/couchbase/gocb.git
```

you can simply perform two simple steps to get started with gerrit:

```
$ git remote add gerrit ssh://${USERNAME}@review.couchbase.org:29418/gocb
$ scp -P 29418 ${USERNAME}@review.couchbase.org:hooks/commit-msg .git/hooks
$ chmod a+x .git/hooks/commit-msg
```

The last change is required for annotating each commit message with a special header known as `Change-Id`.
This allows Gerrit to group together different revisions of the same patch.

#### Pushing a changeset

Now that you have your change and a gerrit account to push to, you need to upload the change for review.
To do so, invoke the following incantation:

```
$ git push gerrit HEAD:refs/for/master
```

Where `gerrit` is the name of the _remote_ added earlier.
You may encounter some errors when pushing.
The most common are:

* "You are not authorized to push to this repository".
  You will get this if your account has not yet been approved.
  Feel free to ask about in gitter.im/couchbase or in the forums for help if blocked.
* "Missing Change-Id". You need to install the `commit-msg` hook as described above.
  Note that even once you do this, you will need to ensure that any prior commits already have this header - this may be
  done by doing an interactive rebase (e.g.  `git rebase -i origin/master` and selecting `reword` for all the commits;
  which will automatically fillin the Change-Id).


Once you've pushed your changeset you can add people to review.
Currently these are:

* Charles Dixon
* Brett Lawson
