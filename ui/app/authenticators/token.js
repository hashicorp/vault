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
  type = 'token';
  displayNamePath = 'display_name';
  tokenPath = 'id';

  async restore(data) {
    // if (data.token) {
    //   return data;
    // }
    throw 'No session stored';
  }

  async login(
    token,
    options = {
      namespace: '',
      backend: 'token',
    }
  ) {
    const url = `/v1/auth/token/lookup-self`;
    const opts = {
      method: 'GET',
      headers: {
        'X-Vault-Token': token,
      },
    };
    if (options.namespace) {
      opts.headers['X-Vault-Namespace'] = options.namespace;
    }
    const result = await fetch(url, opts);
    const body = await result.json();
    if (result.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return this.persistedAuthData(body.data, options);
  }

  persistedAuthData(payload, options) {
    const { entity_id, policies, renewable, namespace_path } = payload;
    const userRootNamespace = this.calculateRootNamespace(options.namespace, namespace_path, options.backend);
    const isRootToken = policies.includes('root');
    const token = payload[this.tokenPath];
    return {
      userRootNamespace,
      isRootToken,
      displayName: payload[this.displayNamePath],
      backend: {
        mountPath: options.backend,
        type: this.type,
      },
      token,
      policies,
      renewable,
      entity_id,
      ...this.calculateExpiration(payload.ttl, payload.lease_duration),
    };
  }

  calculateExpiration(payloadTtl, lease_duration) {
    const now = Date.now();
    const ttl = payloadTtl || lease_duration;
    const tokenExpirationEpoch = now + ttl * 1e3;
    return {
      ttl,
      tokenExpirationEpoch,
    };
  }

  calculateRootNamespace(currentNamespace, namespace_path, backend) {
    // here we prefer namespace_path if its defined,
    // else we look and see if there's already a namespace saved
    // and then finally we'll use the current query param if the others
    // haven't set a value yet
    // all of the typeof checks are necessary because the root namespace is ''
    let userRootNamespace = namespace_path && namespace_path.replace(/\/$/, '');
    // if we're logging in with token and there's no namespace_path, we can assume
    // that the token belongs to the root namespace
    if (backend === 'token' && !userRootNamespace) {
      userRootNamespace = '';
    }
    if (typeof userRootNamespace === 'undefined') {
      // TODO: this doesn't make any sense
      if (this.authData) {
        userRootNamespace = this.authData.userRootNamespace;
      }
    }
    if (typeof userRootNamespace === 'undefined') {
      userRootNamespace = currentNamespace;
    }
    return userRootNamespace;
  }

  async authenticate(token) {
    return this.login(token);
  }
}
