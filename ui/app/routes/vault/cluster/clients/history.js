import Route from '@ember/routing/route';
import { isSameMonth } from 'date-fns';
import RSVP from 'rsvp';
import getStorage from 'vault/lib/token-storage';

const INPUTTED_START_DATE = 'vault:ui-inputted-start-date';
export default class HistoryRoute extends Route {
  async getActivity(start_time) {
    if (isSameMonth(new Date(start_time), new Date())) {
      // triggers empty state to manually enter date if license begins in current month
      return { isLicenseDateError: true };
    }
    // on init ONLY make network request if we have a start_time
    return start_time ? await this.store.queryRecord('clients/activity', { start_time }) : {};
  }

  async getLicenseStartTime() {
    try {
      let license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || getStorage().getItem(INPUTTED_START_DATE) || null;
    } catch (e) {
      // return null so user can input date manually
      // if already inputted manually, will be in localStorage
      return getStorage().getItem(INPUTTED_START_DATE) || null;
    }
  }

  parseRFC3339(timestamp) {
    // convert '2021-03-21T00:00:00Z' --> ['2021', 2] (e.g. 2021 March, month is zero indexed)
    if (Array.isArray(timestamp)) {
      // return if already formatted correctly
      return timestamp;
    }
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
