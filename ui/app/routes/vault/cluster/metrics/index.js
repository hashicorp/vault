import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  model() {
    let entitiesModel = this.store.queryRecord('metrics/entity', {}).then(response => {
      return response.totalEntities;
    });

    let httpsRequestsModel = this.store.queryRecord('metrics/http-requests', {}).then(response => {
      let reverseArray = response.counters.reverse();
      return reverseArray[0].total;
    });
    // ARG TODO: more efficient way to do this.... calling twice, maybe in serializer return the data. or in template?
    let httpsBarChartModel = this.store.queryRecord('metrics/http-requests', {}).then(response => {
      return response.counters;
    });

    let tokenModel = this.store.queryRecord('metrics/token', {}).then(response => {
      return response.totalTokens;
    });

    return hash({
      entitiesTotal: entitiesModel,
      httpsRequestTotal: httpsRequestsModel,
      httpsRequestBarChartData: httpsBarChartModel,
      tokenTotal: tokenModel,
    });
  },
});
