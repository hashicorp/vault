/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

const AUTH = 'vault.cluster.auth';
const PROVIDER = 'vault.cluster.oidc-provider';
const NS_PROVIDER = 'vault.cluster.oidc-provider-ns';

export default class VaultClusterOidcProviderRoute extends Route {
  @service auth;
  @service router;

  get win() {
    return this.window || window;
  }

  _redirect(url, params) {
    if (!url) return;
    const redir = this._buildUrl(url, params);
    if (Ember.testing) {
      return redir;
    }
    this.win.location.replace(redir);
  }

  beforeModel(transition) {
    const currentToken = this.auth.get('currentTokenName');
    const qp = transition.to.queryParams;
    // remove redirect_to if carried over from auth
    qp.redirect_to = null;
    if (!currentToken && 'none' === qp.prompt?.toLowerCase()) {
      this._redirect(qp.redirect_uri, {
        state: qp.state,
        error: 'login_required',
      });
    } else if (!currentToken || 'login' === qp.prompt?.toLowerCase()) {
      const logout = !!currentToken;
      if ('login' === qp.prompt?.toLowerCase()) {
        // need to remove before redirect to avoid infinite loop
        qp.prompt = null;
      }
      return this._redirectToAuth({
        ...transition.to.params,
        qp,
        logout,
      });
    }
  }

  _redirectToAuth({ provider_name, namespace = null, qp, logout = false }) {
    const { cluster_name } = this.paramsFor('vault.cluster');
    let url = namespace
      ? this.router.urlFor(NS_PROVIDER, cluster_name, namespace, provider_name, { queryParams: qp })
      : this.router.urlFor(PROVIDER, cluster_name, provider_name, { queryParams: qp });
    // This is terrible, I'm sorry
    // Need to do this because transitionTo (as used in auth-form) expects url without
    // rootURL /ui/ at the beginning, but urlFor builds it in. We can't use currentRoute
    // because it hasn't transitioned yet
    url = url.replace(/^(\/?ui)/, '');
    if (logout) {
      this.auth.deleteCurrentToken();
    }
    // o param can be anything, as long as it's present the auth page will change
    const queryParams = {
      redirect_to: url,
      o: provider_name,
    };
    if (namespace) {
      queryParams.namespace = namespace;
    }
    return this.transitionTo(AUTH, cluster_name, { queryParams });
  }

  _buildUrl(urlString, params) {
    try {
      const url = new URL(urlString);
      Object.keys(params).forEach((key) => {
        if (params[key]) {
          url.searchParams.append(key, params[key]);
        }
      });
      return url;
    } catch (e) {
      console.debug('DEBUG: parsing url failed for', urlString); // eslint-disable-line
      throw new Error('Invalid URL');
    }
  }

  _handleSuccess(response, baseUrl, state) {
    const { code } = response;
    const redirectUrl = this._buildUrl(baseUrl, { code, state });
    if (!Ember.testing) {
      this.win.location.replace(redirectUrl);
    }
    return { redirectUrl };
  }
  _handleError(errorResp, baseUrl) {
    const redirectUrl = this._buildUrl(baseUrl, { ...errorResp });
    if (!Ember.testing) {
      this.win.location.replace(redirectUrl);
    }
    return { redirectUrl };
  }

  /**
   * Method for getting the parameters from the route. Allows for namespace to be defined on extended route oidc-provider-ns
   * @param {object} params object passed into the model method
   * @returns object with provider_name (string), qp (object of query params), decodedRedirect (string, FQDN)
   */
  _getInfoFromParams(params) {
    const { provider_name, namespace, ...qp } = params;
    const decodedRedirect = decodeURI(qp.redirect_uri);
    return {
      provider_name,
      qp,
      decodedRedirect,
      namespace,
    };
  }

  async model(params) {
    const modelInfo = this._getInfoFromParams(params);
    const { qp, decodedRedirect, ...routeParams } = modelInfo;
    const endpoint = this._buildUrl(
      `${this.win.origin}/v1/identity/oidc/provider/${routeParams.provider_name}/authorize`,
      qp
    );
    if (!qp.redirect_uri) {
      throw new Error('Missing required query params');
    }
    try {
      const response = await this.auth.ajax(endpoint, 'GET', { namespace: routeParams.namespace });
      if ('consent' === qp.prompt?.toLowerCase()) {
        return {
          consent: {
            code: response.code,
            redirect: decodedRedirect,
            state: qp.state,
          },
        };
      }
      return this._handleSuccess(response, decodedRedirect, qp.state);
    } catch (errorRes) {
      const resp = await errorRes.json();
      const code = resp.error;
      if (code === 'max_age_violation' || resp?.errors?.includes('permission denied')) {
        this._redirectToAuth({ ...routeParams, qp, logout: true });
      } else if (code === 'invalid_redirect_uri') {
        return {
          error: {
            title: 'Redirect URI mismatch',
            message:
              'The provided redirect_uri is not in the list of allowed redirect URIs. Please make sure you are sending a valid redirect URI from your application.',
          },
        };
      } else if (code === 'invalid_client_id') {
        return {
          error: {
            title: 'Invalid client ID',
            message: 'Your client ID is invalid. Please update your configuration and try again.',
          },
        };
      } else {
        return this._handleError(resp, decodedRedirect);
      }
    }
  }
}
