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

TODO:

- Move renew/revoke methods back inside expiry (user-menu.hbs)
- Clean up auth service
