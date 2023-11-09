import VaultAuthenticator from './vault-authenticator';

export default class TokenAuthenticator extends VaultAuthenticator {
  type = 'token';
  displayNamePath = 'display_name';
  tokenPath = 'id';

  async login(
    { token },
    options = {
      namespace: '',
      backend: 'token',
    }
  ) {
    const url = `/v1/auth/token/lookup-self`;
    const opts = {
      method: 'GET',
      headers: this.getTokenHeader(token, options.namespace),
    };
    if (options.namespace) {
      opts.headers['X-Vault-Namespace'] = options.namespace;
    }
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return body;
  }
}
