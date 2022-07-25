import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';
import { formatRFC3339 } from 'date-fns';
import { parseDateString } from 'vault/helpers/parse-date-string';

const getActivityParams = ({ tab, start, end }) => {
  // Expects MM-yyyy format
  // TODO: minStart, maxEnd
  let params = {};
  if (tab === 'current') {
    params.tab = tab;
  } else if (tab === 'history') {
    if (start) {
      let startDate = parseDateString(start);
      if (startDate) {
        params.start_time = formatRFC3339(startDate);
      }
    }
    if (end) {
      let endDate = parseDateString(end);
      if (endDate) {
        params.end_time = formatRFC3339(endDate);
      }
    }
  }
  return params;
};

export default Route.extend(ClusterRoute, {
  queryParams: {
    tab: {
      refreshModel: true,
    },
    start: {
      refreshModel: true,
    },
    end: {
      refreshModel: true,
    },
  },

  model(params) {
    let config = this.store.queryRecord('clients/config', {}).catch(e => {
      console.debug(e);
      // swallowing error so activity can show if no config permissions
      return {};
    });
    const activityParams = getActivityParams(params);
    let activity = this.store.queryRecord('clients/activity', activityParams);

    return hash({
      queryStart: params.start,
      queryEnd: params.end,
      activity,
      config,
    });
  },

  actions: {
    loading(transition) {
      // eslint-disable-next-line ember/no-controller-access-in-routes
      let controller = this.controllerFor('vault.cluster.clients.index');
      controller.set('currentlyLoading', true);
      transition.promise.finally(function() {
        controller.set('currentlyLoading', false);
      });
    },
  },
});
