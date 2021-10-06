import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterIdentityOidcProviderRoute extends Route {
  @service auth;
  @service router;

  _redirect(url, params) {
    let redir = this._buildUrl(url, params);
    window.location.replace(redir);
  }

  beforeModel(transition) {
    // const baseUrl = window.location.origin;
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
        console.log('removed prompt');
      }
      let url = this.router.urlFor(transition.to.name, transition.to.params, { queryParams: qp });
      return this.transitionTo('vault.cluster.auth', { queryParams: { redirect_to: url } });
    }
  }
  // beforeModel(transition) {
  //   console.log('beforeModel');
  //   let redirect = this.router.urlFor('vault.cluster.identity.oidc-provider', { state: 'haha' });

  //   console.log({ redirect });
  // this.transitionTo('vault.cluster.auth', { queryParams: { redirect_to: 'foobar' } });
  // let qp = transition.to.queryParams;
  // let currentToken = this.auth.currentTokenName;
  // console.log('BEFORE MODEL');
  // console.log({ qp });
  // console.log({ currentToken });
  // let redirect = this.router.currentURL;
  // if (transition.targetName === this.routeName) {
  //   return this.replaceWith('vault.cluster.auth', { redirect_to: redirect });
  // }

  // if (!currentToken && 'none' === qp.prompt?.toLowerCase()) {
  //   let url = this._buildUrl(qp.redirect_uri, {
  //     state: qp.state,
  //     error: 'login_required',
  //   });
  //   window.location.replace(url);
  // } else if (!currentToken) {
  //   // this.transitionTo('vault.cluster.auth', { redirect_to: redirect });
  //   return this.replaceWith('vault.cluster.identity.oidc-provider', { ...qp });

  //   // this.router.transitionTo(currentRoute);
  //   // let redirect_to = `${window.origin}/ui/vault`;
  //   // window.location.assign(
  //   //   this._buildUrl(`${window.origin}/ui/vault/auth`, {
  //   //     ...qp,
  //   //   })
  //   // );
  // } else if ('login' === qp.prompt?.toLowerCase()) {
  //   // TODO: clear token, redirect to login, redirect_to current route without prompt param
  //   this.transitionTo('vault.cluster.auth', { redirect_to: redirect });
  // }
  // }

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
      console.log('ERROR parsing url', urlString);
      throw new Error('Invalid URL');
    }
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
    console.log('oidc provider name', params);
    let { oidc_name, ...qp } = params;
    console.log('oidc_name', oidc_name);
    let decodedRedirect = decodeURI(qp.redirect_uri);
    let url = this._buildUrl(`${window.origin}/v1/identity/oidc/provider/${oidc_name}/authorize`, qp);
    try {
      // let test = this.router.urlFor('vault.cluster.auth');
      // console.log({ test });
      const response = await this.auth.ajax(url, 'GET', {});
      console.log('CONSENT>>>>>>', qp.prompt);
      if ('consent' === qp.prompt?.toLowerCase()) {
        console.log('SHOULD SHOW CONSENT FORM');
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
      console.error({ errorRes });
      let err = await errorRes.json();
      let code = err.error;
      if (code === 'max_age_violation') {
        // TODO: clear token, redirect to login
        this.auth.deleteCurrentToken();
        // let url = this.router.urlFor(transition.to.name, transition.to.params, { queryParams: qp });
        // return this.transitionTo('vault.cluster.auth', { queryParams: { redirect_to: url } });
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
