import Route from '@ember/routing/route';

export default class VaultClusterOidcProviderRoute extends Route {
  beforeModel(transition) {
    console.log(transition);
  }
  async model(params) {
    console.log({ params });
    let url = new URL(`${window.origin}/v1/identity/oidc/provider/name/authorize`);
    Object.keys(params).forEach(key => {
      if (params[key]) {
        url.searchParams.append(key, params[key]);
      }
    });
    console.log(url, 'URL');
    const result = await fetch('/v1/identity/oidc/provider/name/authorize', {
      method: 'GET',
    });
    // const result = await fetch('v1/identity/oidc/provider/name/authorize', {
    //   method: 'GET',
    // });
    console.log({ result });
    return {};
  }
}
