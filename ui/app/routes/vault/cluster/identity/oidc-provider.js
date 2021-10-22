import Ember from 'ember';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

const AUTH = 'vault.cluster.auth';
const PROVIDER = 'vault.cluster.identity.oidc-provider';

export default class VaultClusterIdentityOidcProviderRoute extends Route {
  @service auth;
  @service router;

  get win() {
    return this.window || window;
  }

  _redirect(url, params) {
    if (!url) return;
    let redir = this._buildUrl(url, params);
    if (Ember.testing) {
      return redir;
    }
    this.win.location.replace(redir);
  }

  beforeModel(transition) {
    const currentToken = this.auth.get('currentTokenName');
    let { redirect_to, ...qp } = transition.to.queryParams;
    console.debug('DEBUG: removing redirect_to', redirect_to);
    if (!currentToken && 'none' === qp.prompt?.toLowerCase()) {
      this._redirect(qp.redirect_uri, {
        state: qp.state,
        error: 'login_required',
      });
    } else if (!currentToken || 'login' === qp.prompt?.toLowerCase()) {
      let shouldLogout = !!currentToken;
      if ('login' === qp.prompt?.toLowerCase()) {
        // need to remove before redirect to avoid infinite loop
        qp.prompt = null;
      }
      return this._redirectToAuth(transition.to.params?.provider_name, qp, shouldLogout);
    }
  }

  _redirectToAuth(provider_name, queryParams, logout = false) {
    let { cluster_name } = this.paramsFor('vault.cluster');
    let url = this.router.urlFor(PROVIDER, cluster_name, provider_name, { queryParams });
    // This is terrible, I'm sorry
    // Need to do this because transitionTo (as used in auth-form) expects url without
    // rootURL /ui/ at the beginning, but urlFor builds it in. We can't use currentRoute
    // because it hasn't transitioned yet
    url = url.replace(/^(\/?ui)/, '');
    if (logout) {
      this.auth.deleteCurrentToken();
    }
    // o param can be anything, as long as it's present the auth page will change
    return this.transitionTo(AUTH, cluster_name, { queryParams: { redirect_to: url, o: provider_name } });
  }

  _buildUrl(urlString, params) {
    try {
      let url = new URL(urlString);
      Object.keys(params).forEach(key => {
        if (params[key]) {
          url.searchParams.append(key, params[key]);
        }
      });
      return url;
    } catch (e) {
      console.debug('DEBUG: parsing url failed for', urlString);
      throw new Error('Invalid URL');
    }
  }

  _handleSuccess(response, baseUrl, state) {
    const { code } = response;
    let redirectUrl = this._buildUrl(baseUrl, { code, state });
    if (Ember.testing) {
      return redirectUrl;
    }
    this.win.location.replace(redirectUrl);
  }
  _handleError(errorResp, baseUrl) {
    let redirectUrl = this._buildUrl(baseUrl, { ...errorResp });
    if (Ember.testing) {
      return redirectUrl;
    }
    this.win.location.replace(redirectUrl);
  }

  async model(params) {
    let { provider_name, ...qp } = params;
    let decodedRedirect = decodeURI(qp.redirect_uri);
    let baseUrl = this.namespace.path
      ? `${this.win.origin}/v1/${this.namespace.path}/identity/oidc/provider/${provider_name}/authorize`
      : `${this.win.origin}/v1/identity/oidc/provider/${provider_name}/authorize`;
    let endpoint = this._buildUrl(baseUrl, qp);
    try {
      const response = await this.auth.ajax(endpoint, 'GET', {});
      if ('consent' === qp.prompt?.toLowerCase()) {
        return {
          consent: {
            code: response.code,
            redirect: decodedRedirect,
            state: qp.state,
          },
        };
      }
      this._handleSuccess(response, decodedRedirect, qp.state);
    } catch (errorRes) {
      let resp = await errorRes.json();
      let code = resp.error;
      if (code === 'max_age_violation' || resp?.errors?.includes('permission denied')) {
        this._redirectToAuth(provider_name, qp, true);
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
        this._handleError(resp, decodedRedirect);
      }
    }
  }
}
