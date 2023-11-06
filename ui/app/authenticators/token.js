import BaseAuthenticator from 'ember-simple-auth/authenticators/base';

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
export default class TokenAuthenticator extends BaseAuthenticator {
  async login(token) {
    const url = `/v1/auth/token/lookup-self`;
    const result = await fetch(url, {
      method: 'POST',
      headers: {
        'X-Vault-Token': token,
      },
    });
    const body = await result.json();
    console.log({ result, body });
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return body;
  }

  async authenticate(token) {
    return this.login(token);
  }
}
