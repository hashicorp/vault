import Route from '@ember/routing/route';
import RSVP from 'rsvp';
export default class OidcClientDetailsRoute extends Route {
  queryParams = {
    tab: {
      refreshModel: true,
    },
  };

  model(params) {
    const { tab } = params;
    const model = this.modelFor('vault.cluster.access.oidc.clients.client');
    if (tab === 'details') {
      return model;
    }
    if (tab === 'providers') {
      return RSVP.hash({
        name: model.name,
        providers: this.store.query('oidc/provider', {
          allowed_client_id: model.clientId,
        }),
      });
    }
  }

  // Reset query params to default since query param values in Ember are "sticky"
  // and the latest query param is preserved,
  resetController(controller, isExiting) {
    if (isExiting) {
      // isExiting is false if only the route's model was changing
      controller.set('tab', 'details');
    }
  }
}
