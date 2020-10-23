# Combined Database Engine
This package is how database plugins interact with Vault.

# Upgrading to Version 5
In Vault 1.6, a new Database interface was created that solved a number of issues with the previous interface:

1. It could not use password policies because the database plugins were responsible for generating passwords.
2. There were significant inconsistencies between functions in the interface.
3. Several functions (`SetCredentials` and `RotateRootCredentials`) were doing the same operation.
4. It had a function that was no longer being used as it had been deprecated in a previous version but never removed.

We highly recommend that you upgrade any version 4 database plugins to version 5 as version 4 is considered deprecated
and support for it will be removed in a future release.


