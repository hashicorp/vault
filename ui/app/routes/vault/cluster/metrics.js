import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  queryParams: {
    start: {
      refreshModel: true,
    },
    end: {
      refreshModel: true,
    },
  },

  model(params) {
    let config = this.store.queryRecord('metrics/config', {});

    let activity = this.store.queryRecord('metrics/activity', params);

    return hash({
      activity,
      config,
    });
  },
});
