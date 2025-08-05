/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { service } from '@ember/service';
import { hash, resolve } from 'rsvp';
import { assert } from '@ember/debug';
import { pluralize } from 'ember-inflector';

import ApplicationAdapter from './application';

const ENDPOINTS = [
  'health',
  'seal-status',
  'tokens',
  'token',
  'seal',
  'unseal',
  'init',
  'capabilities-self',
  'license',
];

const REPLICATION_ENDPOINTS = {
  reindex: 'reindex',
  recover: 'recover',
  status: 'status',

  primary: ['enable', 'disable', 'demote', 'secondary-token', 'revoke-secondary'],

  secondary: ['enable', 'disable', 'promote', 'update-primary'],
};

const REPLICATION_MODES = ['dr', 'performance'];
export default ApplicationAdapter.extend({
  version: service(),
  namespaceService: service('namespace'),
  shouldBackgroundReloadRecord() {
    return true;
  },

  findRecord(store, type, id, snapshot) {
    const fetches = {
      health: this.health(),
      sealStatus: this.sealStatus().catch((e) => e),
    };
    if (this.version.isEnterprise && this.namespaceService.inRootNamespace) {
      fetches.replicationStatus = this.replicationStatus().catch((e) => e);
    }
    return hash(fetches).then(({ health, sealStatus, replicationStatus }) => {
      let ret = {
        id,
        name: snapshot.attr('name'),
      };
      ret = Object.assign(ret, health);
      if (sealStatus instanceof AdapterError === false) {
        ret = Object.assign(ret, { nodes: [sealStatus] });
      }
      if (replicationStatus && replicationStatus instanceof AdapterError === false) {
        ret = Object.assign(ret, replicationStatus.data);
      } else if (
        replicationStatus instanceof AdapterError &&
        replicationStatus?.errors.find((err) => err === 'disabled path')
      ) {
        // set redacted if result is an error which only happens when redacted
        ret = Object.assign(ret, { replication_redacted: true });
      }
      return resolve(ret);
    });
  },

  pathForType(type) {
    return type === 'cluster' ? 'clusters' : pluralize(type);
  },

  health() {
    return this.ajax(this.urlFor('health'), 'GET', {
      data: {
        standbycode: 200,
        sealedcode: 200,
        uninitcode: 200,
        drsecondarycode: 200,
        performancestandbycode: 200,
      },
      unauthenticated: true,
    }).catch(() => {
      // sys/health will only fail when chroot set
      // because it's allowed in root namespace only and
      // configured to return a 200 response in other fail scenarios
      return { has_chroot_namespace: true };
    });
  },

  features() {
    return this.ajax(`${this.urlFor('license')}/features`, 'GET', {
      unauthenticated: true,
    });
  },

  sealStatus(unauthenticated = true) {
    return this.ajax(this.urlFor('seal-status'), 'GET', { unauthenticated });
  },

  seal() {
    return this.ajax(this.urlFor('seal'), 'PUT');
  },

  unseal(data) {
    return this.ajax(this.urlFor('unseal'), 'PUT', {
      data,
      unauthenticated: true,
    });
  },

  initCluster(data) {
    return this.ajax(this.urlFor('init'), 'PUT', {
      data,
      unauthenticated: true,
    });
  },

  authenticate({ backend, data }) {
    const { role, jwt, token, password, username, path, nonce } = data;
    const url = this.urlForAuth(backend, username, path);
    const verb = backend === 'token' ? 'GET' : 'POST';
    const options = {
      unauthenticated: true,
    };
    if (backend === 'token') {
      options.headers = {
        'X-Vault-Token': token,
      };
    } else if (backend === 'jwt' || backend === 'oidc') {
      options.data = { role, jwt };
    } else if (backend === 'okta') {
      options.data = { password, nonce };
    } else {
      options.data = token ? { token, password } : { password };
    }

    return this.ajax(url, verb, options);
  },

  mfaValidate({ mfa_request_id, mfa_constraints }) {
    const options = {
      data: {
        mfa_request_id,
        mfa_payload: mfa_constraints.reduce((obj, { selectedMethod, passcode }) => {
          let payload = [];
          if (passcode) {
            // duo requires passcode= prepended to the actual passcode
            // this isn't a great UX so we add it behind the scenes to fulfill the requirement
            // check if user added passcode= to avoid duplication
            payload =
              selectedMethod.type === 'duo' && !passcode.includes('passcode=')
                ? [`passcode=${passcode}`]
                : [passcode];
          }
          obj[selectedMethod.id] = payload;
          return obj;
        }, {}),
      },
    };
    return this.ajax('/v1/sys/mfa/validate', 'POST', options);
  },

  urlFor(endpoint) {
    if (!ENDPOINTS.includes(endpoint)) {
      throw new Error(
        `Calls to a ${endpoint} endpoint are not currently allowed in the vault cluster adapter`
      );
    }
    return `${this.buildURL()}/${endpoint}`;
  },

  urlForAuth(type, username, path) {
    const authBackend = type.toLowerCase();
    const authURLs = {
      github: 'login',
      jwt: 'login',
      oidc: 'login',
      userpass: `login/${encodeURIComponent(username)}`,
      ldap: `login/${encodeURIComponent(username)}`,
      okta: `login/${encodeURIComponent(username)}`,
      radius: `login/${encodeURIComponent(username)}`,
      token: 'lookup-self',
    };
    const urlSuffix = authURLs[authBackend];
    const urlPrefix = path && authBackend !== 'token' ? path : authBackend;
    if (!urlSuffix) {
      throw new Error(`There is no auth url for ${type}.`);
    }
    return `/v1/auth/${urlPrefix}/${urlSuffix}`;
  },

  urlForReplication(replicationMode, clusterMode, endpoint) {
    let suffix;
    const errString = `Calls to replication ${endpoint} endpoint are not currently allowed in the vault cluster adapater`;
    if (clusterMode) {
      assert(errString, REPLICATION_ENDPOINTS[clusterMode].includes(endpoint));
      suffix = `${replicationMode}/${clusterMode}/${endpoint}`;
    } else {
      assert(errString, REPLICATION_ENDPOINTS[endpoint]);
      suffix = `${endpoint}`;
    }
    return `${this.buildURL()}/replication/${suffix}`;
  },

  replicationStatus() {
    return this.ajax(`${this.buildURL()}/replication/status`, 'GET', { unauthenticated: true });
  },

  replicationDrPromote(data, options) {
    const verb = options && options.checkStatus ? 'GET' : 'PUT';
    return this.ajax(`${this.buildURL()}/replication/dr/secondary/promote`, verb, {
      data,
      unauthenticated: true,
    });
  },

  generateDrOperationToken(data, options) {
    let verb = 'POST';
    let url = `${this.buildURL()}/replication/dr/secondary/generate-operation-token/`;
    if (options?.cancel) {
      verb = 'DELETE';
      url += 'attempt';
    } else if (options?.checkStatus) {
      verb = 'GET';
      url += 'attempt';
    } else if (data?.pgp_key || data?.attempt) {
      url += 'attempt';
    } else {
      // progress the operation
      url += 'update';
    }
    return this.ajax(url, verb, {
      data,
      unauthenticated: true,
    });
  },

  replicationAction(action, replicationMode, clusterMode, data) {
    assert(
      `${replicationMode} is an unsupported replication mode.`,
      replicationMode && REPLICATION_MODES.includes(replicationMode)
    );

    const url =
      action === 'recover' || action === 'reindex'
        ? this.urlForReplication(replicationMode, null, action)
        : this.urlForReplication(replicationMode, clusterMode, action);

    return this.ajax(url, 'POST', { data });
  },
});
