import Route from '@ember/routing/route';
import RSVP from 'rsvp';

export default class HistoryRoute extends Route {
  async getActivity(start_time) {
    try {
      // on init ONLY make network request if we have a start time from the license
      // otherwise user needs to manually input
      return start_time
        ? await this.store.queryRecord('clients/activity', { start_time })
        : { endTime: null };
    } catch (e) {
      return e;
    }
  }

  async getLicenseStartTime() {
    try {
      let license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || null;
    } catch (e) {
      // if error due to permission denied, return null so user can input date manually
      return null;
    }
  }

  parseRFC3339(timestamp) {
    // convert '2021-03-21T00:00:00Z' --> ['2021', 2] (e.g. 2021 March, month is zero indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  }

  async model() {
    let parentModel = this.modelFor('vault.cluster.clients');
    let licenseStart = await this.getLicenseStartTime();
    let activity = await this.getActivity(licenseStart);

    return RSVP.hash({
      config: parentModel.config,
      activity,
      startTimeFromLicense: this.parseRFC3339(licenseStart),
      endTimeFromResponse: this.parseRFC3339(activity?.endTime),
      versionHistory: parentModel.versionHistory,
    });
  }
}
