import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';

export default class HistoryRoute extends Route {
  async getLicense() {
    try {
      return this.store.queryRecord('license', {});
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  }

  async getActivity(start_time) {
    try {
      return this.store.queryRecord('clients/activity', { start_time });
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  }

  parseRFC3339(timestamp) {
    // convert '2021-03-21T00:00:00Z' --> ['2021', 2] (e.g. 2021 March, month is zero indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  }

  async model() {
    let license = await this.getLicense();
    let activity = await this.getActivity(license.startTime);

    return RSVP.hash({
      config: this.store.queryRecord('clients/config', {}),
      activity,
      startTimeFromLicense: this.parseRFC3339(license.startTime),
      endTimeFromResponse: activity ? this.parseRFC3339(activity.endTime) : null,
    });
  }

  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor('vault.cluster.clients.history');
    controller.set('currentlyLoading', true);
    transition.promise.finally(function () {
      controller.set('currentlyLoading', false);
    });
  }
}
