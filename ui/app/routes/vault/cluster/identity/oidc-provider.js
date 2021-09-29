import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterIdentityOidcProviderRoute extends Route {
  @service auth;

  beforeModel(transition) {
    console.log(transition);
  }

  _buildUrl(urlString, params) {
    let url = new URL(urlString);
    Object.keys(params).forEach(key => {
      if (params[key]) {
        url.searchParams.append(key, params[key]);
      }
    });
    return url;
  }

  _handleSuccess(response, baseUrl, state) {
    const { code } = response;
    let redirectUrl = this._buildUrl(baseUrl, { code, state });
    console.log('REDIRECT TO (SUCCESS)', redirectUrl);
    // window.location.replace(redirectUrl);
  }
  _handleError(response, baseUrl) {
    const { code } = response;
    let redirectUrl = this._buildUrl(baseUrl, { error: code });
    console.log('REDIRECT TO (ERROR)', redirectUrl);
    // window.location.replace(redirectUrl);
  }

  async model(params) {
    let currentToken = this.auth.currentTokenName;
    let { oidc_name, ...qp } = params;

    if (!currentToken && 'none' === qp.prompt?.toLowerCase()) {
      // TODO: show error (return error on model?)
    } else if (!currentToken) {
      // TODO: redirect to login
    } else if ('login' === qp.prompt?.toLowerCase()) {
      // TODO: clear token, redirect to login, redirect_to current route without prompt param
    }

    let decodedRedirect = decodeURI(qp.redirect_uri);
    let url = this._buildUrl(`${window.origin}/v1/identity/oidc/provider/${oidc_name}/authorize`, qp);
    try {
      const response = await this.auth.ajax(url, 'GET', {});
      console.log(response);
      if ('consent' === qp.prompt?.toLowerCase()) {
        // TODO: Show consent before redirect (code val must be passed to action)
        return {
          consent: {
            code: response.code,
            redirect: decodedRedirect,
            state: qp.state,
          },
        };
      } else {
        this._handleSuccess(response, decodedRedirect, qp.state);
      }
    } catch (errorRes) {
      let err = await errorRes.json();
      let code = err.error;
      if (code === 'max_age_violation') {
        // TODO: clear token, redirect to login
      } else if (code === 'invalid_redirect_uri') {
        // TODO: show error (return on model)
      } else if (code === 'invalid_client_id') {
        // TODO: show error (return on model)
      } else {
        // TODO: return error as param to callback
      }
      return err;
    }
  }
}
