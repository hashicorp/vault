import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  store: service(),
  beforeModel() {
    let { auth_path, state, code } = this.paramsFor(this.routeName);
    let adapter = this.store.adapterFor('auth-method');

    return adapter.exchangeOIDC(auth_path, state, code).then(resp => {
      let token = resp.wrapped_token;
      return this.transitionTo('vault.cluster.auth', { queryParams: { wrapped_token: token } });
    });
  },
});
