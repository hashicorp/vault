import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  /* Asynchronous example
  model: function(params) {
    var reviewPromise = this.store.findRecord('review', params.id);
    return Ember.RSVP.hash({
      review: reviewPromise,
      user: reviewPromise.then(review => {
        return this.store.findRecord('user', review.get('userId'));
      })
    });
  }
  */
  model() {
    let totalEntities = this.store.queryRecord('metrics/entity', {}).then(response => {
      return response.entities.total;
    });

    let httpsRequests = this.store.queryRecord('metrics/http-requests', {}).then(response => {
      let reverseArray = response.counters.reverse();
      return reverseArray;
    });

    let totalTokens = this.store.queryRecord('metrics/token', {}).then(response => {
      return response.service_tokens.total;
    });

    let activity = this.store
      .queryRecord('metrics/activity', {})
      .then(response => {
        console.log({ response });
        return response.data;
      })
      .catch(e => {
        return {
          active_clients: 1000,
          unique_entities: 635,
          direct_tokens: 465,
        };
      });

    return hash({
      activity,
      totalEntities,
      httpsRequests,
      totalTokens,
    });
  },
});
