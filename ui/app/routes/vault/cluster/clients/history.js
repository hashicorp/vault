import Route from '@ember/routing/route';
import { isSameMonth } from 'date-fns';
import RSVP from 'rsvp';
import getStorage from 'vault/lib/token-storage';
import { parseRFC3339 } from 'core/utils/date-formatters';

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

  async model() {
    let parentModel = this.modelFor('vault.cluster.clients');
    let licenseStart = await this.getLicenseStartTime();
    let activity = await this.getActivity(licenseStart);

    return RSVP.hash({
      config: parentModel.config,
      activity,
      startTimeFromLicense: parseRFC3339(licenseStart),
      endTimeFromResponse: parseRFC3339(activity?.endTime),
      versionHistory: parentModel.versionHistory,
    });
  }
}
