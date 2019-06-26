# Git hooks

The files in .hooks/git/ are the raw git hooks that get executed by git.
They are installed by the Makefile in the root of this repo.

Those git hooks call the files in .hooks/ prefixed with the hook name.
E.g. `.hooks/git/pre-commit` calls all the files prefixed with `.hooks/pre-commit-*`

