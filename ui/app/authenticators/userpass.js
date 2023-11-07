import VaultAuthenticator from './vault-authenticator';

export default class UserpassAuthenticator extends VaultAuthenticator {
  type = 'userpass';
  displayNamePath = 'metadata.username';
  tokenPath = 'client_token';

  async login(token, options) {
    const url = `/v1/auth/userpass/login/${encodeURIComponent(options.username)}`;
    const opts = {
      method: 'POST',
      headers: this.getTokenHeader(token, options.namespace),
      body: JSON.stringify({ password: token }),
    };
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return this.persistedAuthData(body.auth, options);
  }
}
