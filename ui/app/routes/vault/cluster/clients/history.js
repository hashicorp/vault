import Route from '@ember/routing/route';
import { isSameMonth } from 'date-fns';
import RSVP from 'rsvp';
import getStorage from 'vault/lib/token-storage';
import { parseRFC3339 } from 'core/utils/date-formatters';
import { inject as service } from '@ember/service';
const INPUTTED_START_DATE = 'vault:ui-inputted-start-date';

export default class HistoryRoute extends Route {
  @service store;

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
      const license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || getStorage().getItem(INPUTTED_START_DATE) || null;
    } catch (e) {
      // return null so user can input date manually
      // if already inputted manually, will be in localStorage
      return getStorage().getItem(INPUTTED_START_DATE) || null;
    }
  }

  async model() {
    const parentModel = this.modelFor('vault.cluster.clients');
    const licenseStart = await this.getLicenseStartTime();
    const activity = await this.getActivity(licenseStart);

    return RSVP.hash({
      config: parentModel.config,
      activity,
      startTimeFromLicense: parseRFC3339(licenseStart),
      endTimeFromResponse: parseRFC3339(activity?.endTime),
      versionHistory: parentModel.versionHistory,
    });
  }
}
