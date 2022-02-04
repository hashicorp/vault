import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';

export default class HistoryRoute extends Route {
  async getLicense() {
    try {
      return await this.store.queryRecord('license', {});
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  }

  async getActivity(start_time) {
    try {
      return await this.store.queryRecord('clients/activity', { start_time });
    } catch (e) {
      // ARG TODO handle
      return e;
    }
  }

  rfc33395ToMonthYear(timestamp) {
    // return ['2021', 2] (e.g. 2021 March, make 0-indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  }

  async model() {
    let license = await this.getLicense(); // get default start_time
    let activity = await this.getActivity(license.startTime); // returns client counts using license start_time.
    let endTimeFromResponse = activity ? this.rfc33395ToMonthYear(activity.endTime) : null;
    let startTimeFromLicense = this.rfc33395ToMonthYear(license.startTime);

    return RSVP.hash({
      config: this.store.queryRecord('clients/config', {}),
      activity,
      startTimeFromLicense,
      endTimeFromResponse,
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
