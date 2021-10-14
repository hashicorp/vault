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
    let redir = this._buildUrl(url, params);
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
      if ('login' === qp.prompt?.toLowerCase()) {
        this.auth.deleteCurrentToken();
        qp.prompt = null;
      }
      let { cluster_name } = this.paramsFor('vault.cluster');
      let url = this.router.urlFor(transition.to.name, cluster_name, transition.to.params, {
        queryParams: qp,
      });
      return this.transitionTo(AUTH, cluster_name, { queryParams: { redirect_to: url } });
    }
  }

  _redirectToAuth(oidcName, queryParams, logout = false) {
    let { cluster_name } = this.paramsFor('vault.cluster');
    let currentRoute = this.router.urlFor(PROVIDER, cluster_name, oidcName, { queryParams });
    if (logout) {
      this.auth.deleteCurrentToken();
    }
    return this.transitionTo(AUTH, cluster_name, { queryParams: { redirect_to: currentRoute } });
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
    this.win.location.replace(redirectUrl);
  }
  _handleError(errorResp, baseUrl) {
    let redirectUrl = this._buildUrl(baseUrl, { ...errorResp });
    this.win.location.replace(redirectUrl);
  }

  async model(params) {
    let { oidc_name, ...qp } = params;
    let decodedRedirect = decodeURI(qp.redirect_uri);
    let url = this._buildUrl(`${this.win.origin}/v1/identity/oidc/provider/${oidc_name}/authorize`, qp);
    try {
      const response = await this.auth.ajax(url, 'GET', {});
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
      if (code === 'max_age_violation') {
        this._redirectToAuth(oidc_name, qp, true);
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
