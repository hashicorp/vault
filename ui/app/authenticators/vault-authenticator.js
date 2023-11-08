import BaseAuthenticator from 'ember-simple-auth/authenticators/base';
import RSVP from 'rsvp';
import { get } from '@ember/object';
import { assert } from '@ember/debug';

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
export default class VaultAuthenticator extends BaseAuthenticator {
  getTokenHeader(token, namespace) {
    const headers = {
      'X-Vault-Token': token,
    };
    if (namespace) {
      headers['X-Vault-Namespace'] = namespace;
    }
    return headers;
  }

  async restore(data) {
    if (data.token) {
      return data;
    }
    throw 'No session stored';
  }

  async authenticate(fields, options) {
    const { renew, ...opts } = options;
    if (renew) {
      const renewed = await this.renewToken(fields.token, opts.namespace);
      // TODO: OIDC renew endpoint doesn't return display_name
      return this.persistedAuthData(renewed, opts, 'client_token');
    }
    return this.login(fields, opts);
  }

  invalidate(authData, options) {
    if (options?.revoke) {
      return this.revokeToken(authData.token, options.namespace);
    }
    return RSVP.resolve();
  }

  /* methods */
  async revokeToken(token, namespace) {
    const url = '/v1/auth/token/revoke-self';
    const headers = this.getTokenHeader(token, namespace);

    const response = await fetch(url, {
      method: 'POST',
      headers,
    });
    if (response.status !== 204) {
      const body = await response.json();
      throw new Error(body.errors.join(', '));
    }
    return;
  }

  async renewToken(token, namespace) {
    const url = '/v1/auth/token/renew-self';
    const headers = this.getTokenHeader(token, namespace);
    const response = await fetch(url, {
      method: 'POST',
      headers,
    });
    const body = await response.json();
    if (response.status !== 200) {
      throw new Error(body.errors.join(', '));
    }
    return body.data || body.auth;
  }

  persistedAuthData(data, options, tokenPath = 'id') {
    if (!options?.backend) {
      assert('persistedAuthData requires options with backend');
    }
    const { entity_id, policies, renewable, namespace_path } = data;
    const userRootNamespace = this.calculateRootNamespace(options.namespace, namespace_path, options.backend);
    const isRootToken = policies.includes('root');
    const token = data[tokenPath];

    const persisted = {
      userRootNamespace,
      isRootToken,
      displayName: get(data, this.displayNamePath),
      backend: {
        mountPath: options.backend,
        type: this.type,
      },
      token,
      policies,
      renewable,
      entity_id,
      ...this.calculateExpiration(data.ttl, data.lease_duration),
    };
    return persisted;
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
    // here we prefer namespace_path from the payload if its defined,
    // else we'll use the current query param if the others
    // haven't set a value yet
    // all of the typeof checks are necessary because the root namespace is ''
    let userRootNamespace = namespace_path && namespace_path.replace(/\/$/, '');
    // if we're logging in with token and there's no namespace_path, we can assume
    // that the token belongs to the root namespace
    if (backend === 'token' && !userRootNamespace) {
      userRootNamespace = '';
    }
    if (typeof userRootNamespace === 'undefined') {
      userRootNamespace = currentNamespace;
    }
    return userRootNamespace;
  }
}
