import VaultAuthenticator from './vault-authenticator';

export default class OidcAuthenticator extends VaultAuthenticator {
  type = 'oidc';
  displayNamePath = 'display_name';
  tokenPath = 'client_token';

  async lookupSelf(
    token,
    options = {
      namespace: '',
      backend: 'oidc',
    }
  ) {
    const url = `/v1/auth/token/lookup-self`;
    const opts = {
      method: 'GET',
      headers: this.getTokenHeader(token, options.namespace),
    };
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return this.persistedAuthData(body.data, options, 'id');
  }

  async login(
    { state, code },
    options = {
      namespace: '',
      backend: 'oidc',
    }
  ) {
    if (!state || !code) {
      throw new Error('missing authorization code parameters');
    }
    const url = `/v1/auth/${encodeURIComponent(options.backend)}/oidc/callback?state=${state}&code=${code}`;
    const opts = {
      method: 'GET',
    };
    if (options.namespace) {
      opts.headers['X-Vault-Namespace'] = options.namespace;
    }
    // Exchange the authorization code for an OIDC ID Token
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    // Finally, lookup the token returned from the oidc endpoint
    // so the UI has all the data it needs
    return this.lookupSelf(body.auth[this.tokenPath], options);
  }
}
