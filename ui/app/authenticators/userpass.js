import VaultAuthenticator from './vault-authenticator';

/*{
  "body": {
    "data": {
      "accessor": "dOtbwJJ1Uz2oB9I21vV0LemT",
      "creation_time": 1699288223,
      "creation_ttl": 0,
      "display_name": "token",
      "entity_id": "",
      "expire_time": null,
      "explicit_max_ttl": 0,
      "id": "root",
      "issue_time": "2023-11-06T10:30:23.328432-06:00",
      "meta": null,
      "num_uses": 0,
      "orphan": true,
      "path": "auth/token/create",
      "policies": [
        "root"
      ],
      "renewable": false,
      "ttl": 0,
      "type": "service"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
  }
}
*/
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
