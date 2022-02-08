import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

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

  async getActivity(start_time) {
    try {
      return this.store.queryRecord('clients/activity', { start_time });
    } catch (e) {
      return e;
    }
  },

  async getVersionHistory() {
    try {
      let arrayOfModels = [];
      let response = await this.store.findAll('clients/version-history'); // returns a class with nested models
      response.forEach((model) => {
        arrayOfModels.push({
          id: model.id,
          perviousVersion: model.previousVersion,
          timestampInstalled: model.timestampInstalled,
        });
      });
      return arrayOfModels;
    } catch (e) {
      return null;
    }
  },

  async getLicense() {
    try {
      return await this.store.queryRecord('license', {});
    } catch (e) {
      return e;
    }
  },

  async getMonthly() {
    try {
      return await this.store.queryRecord('clients/monthly', {});
    } catch (e) {
      return e;
    }
  },

  rfc33395ToMonthYear(timestamp) {
    // return ['2021', 2] (e.g. 2021 March, make 0-indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  },

  async model() {
    let config = this.store.queryRecord('clients/config', {}).catch((e) => {
      console.debug(e);
      // swallowing error so activity can show if no config permissions
      return {};
    });

    let license = await this.getLicense(); // get default start_time
    let activity = await this.getActivity(license.startTime); // returns client counts using license start_time.
    let monthly = await this.getMonthly(); // returns the partial month endpoint
    let versionHistory = await this.getVersionHistory();
    let endTimeFromResponse = activity ? this.rfc33395ToMonthYear(activity.endTime) : null;
    let startTimeFromLicense = this.rfc33395ToMonthYear(license.startTime);

    return hash({
      // ARG TODO will remove "hash" once remove "activity," which currently relies on it.
      activity,
      monthly,
      config,
      endTimeFromResponse,
      startTimeFromLicense,
      versionHistory,
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
