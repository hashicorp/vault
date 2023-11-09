import VaultAuthenticator from './vault-authenticator';

export default class UserpassAuthenticator extends VaultAuthenticator {
  type = 'userpass';
  displayNamePath = 'metadata.username';
  tokenPath = 'client_token';

  async login({ username, password }, options) {
    const url = `/v1/auth/${options.backend}/login/${encodeURIComponent(username)}`;
    const opts = {
      method: 'POST',
      headers: this.getTokenHeader(password, options.namespace),
      body: JSON.stringify({ password }),
    };
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return body;
  }
}
