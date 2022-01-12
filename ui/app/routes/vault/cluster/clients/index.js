import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';
import { getTime } from 'date-fns';
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
        // TODO: Replace with formatRFC3339 when date-fns is updated
        // converts to milliseconds, divide by 1000 to get epoch
        params.start_time = getTime(startDate) / 1000;
      }
    }
    if (end) {
      let endDate = parseDateString(end);
      if (endDate) {
        // TODO: Replace with formatRFC3339 when date-fns is updated
        params.end_time = getTime(endDate) / 1000;
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

  async getLicense() {
    try {
      return await this.store.queryRecord('license', {});
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  },

  async getNewInitActivity(start_time) {
    try {
      return await this.store.queryRecord('clients/newInitActivity', { start_time });
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  },

  // ARG POC
  async model(params) {
    let config = this.store.queryRecord('clients/config', {}).catch((e) => {
      console.debug(e);
      // swallowing error so activity can show if no config permissions
      return {};
    });
    const activityParams = getActivityParams(params);
    let activity = this.store.queryRecord('clients/activity', activityParams); // ARG TODO should async await and try catch
    let license = await this.getLicense();
    let newInitActivity = await this.getNewInitActivity(license.startTime);
    // ARG TODO this is the new pattern, while keeping the old. On index load, hit sys/license/status to check for a startTime and return it.

    return hash({
      queryStart: params.start, // ARG TODO will remove
      queryEnd: params.end, // ARG TODO will remove
      activity,
      newInitActivity,
      config,
      startTime: license.startTime,
    });
  },

  actions: {
    loading(transition) {
      // eslint-disable-next-line ember/no-controller-access-in-routes
      let controller = this.controllerFor('vault.cluster.clients.index');
      controller.set('currentlyLoading', true);
      transition.promise.finally(function () {
        controller.set('currentlyLoading', false);
      });
    },
  },
});
