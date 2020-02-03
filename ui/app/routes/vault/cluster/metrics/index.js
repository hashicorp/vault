import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  model() {
    let totalEntities = this.store.queryRecord('metrics/entity', {}).then(response => {
      return response.entities.total;
    });

    let totalHttpsRequest = this.store.queryRecord('metrics/http-requests', {}).then(response => {
      let reverseArray = response.counters.reverse();
      return reverseArray[0].total;
    });

    let totalTokens = this.store.queryRecord('metrics/token', {}).then(response => {
      console.log(response, 'tokens');
      return response.service_tokens.total;
    });

    return hash({
      totalEntities,
      totalHttpsRequest,
      totalTokens,
    });
  },
});
