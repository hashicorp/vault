# Add ember-simple-auth to Vault

## Resources

- docs https://ember-simple-auth.com/
- readme https://github.com/mainmatter/ember-simple-auth/tree/5.0.0

# Plan

- Install ember-simple-auth
- create new auth component for the form, and get it working with token login
- enable session store (think about how to make it configurable?)
- Move files to V2 Addon?

# Steps

1. `ember install ember-simple-auth`
2. Made AuthV2 component and got it working for token
3. Made token authenticator which works for login and error handling
4. Persisted auth data for authenticated users
5. Replaced isActiveSession from the cluster and application template
6. Added session.invaildate to the logout route
7. Added logic for revoke and renew
8. Add base authenticator for common/shared methods
9. Created auth-form components by type and decorator
10. Got OIDC login working
11. Worked on fixing up renew token flow (it may have a bug on OIDC)
12. Got userpass working
13. Cleaned up auth service
14. Trigger onSuccess and onUpdate side effects for query param updates

TODO:

- Nice-to-have show tabs for auth methods that configured that option (Kianna)
- Handle wrapped token (Chelsea)
- MFA (Chelsea)
- Okta challenge (Kianna)
-
