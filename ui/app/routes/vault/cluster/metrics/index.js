import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  model() {
    let entitiesModel = this.store.queryRecord('entities', {}).then(response => {
      return response.counters.entities.total;
    });

    let httpsModel = this.store.queryRecord('http-requests', {}).then(response => {
      let reverseArray = response.counters.reverse();
      return reverseArray[0].total;
    });

    let tokenModel = this.store.queryRecord('tokens', {}).then(response => {
      return response.counters.service_tokens.total;
    });

    return hash({
      entitiesTotal: entitiesModel,
      httpsRequestTotal: httpsModel,
      tokenTotal: tokenModel,
    });
  },
});
