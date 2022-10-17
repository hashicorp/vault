import Route from '@ember/routing/route';
import getStorage from 'vault/lib/token-storage';

// TODO CMB: change class and file name to dashboard
export default class HistoryRoute extends Route {
  currentDate = new Date().toISOString();

  async getActivity(start_time) {
    // on init ONLY make network request if we have a start_time
    return start_time
      ? await this.store.queryRecord('clients/activity', { start_time, end_time: this.currentDate })
      : {};
  }

  async getLicenseStartTime() {
    try {
      let license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || getStorage().getItem('vault:ui-inputted-start-date') || null;
    } catch (e) {
      // return null so user can input date manually
      // if already inputted manually, will be in localStorage
      return getStorage().getItem('vault:ui-inputted-start-date') || null;
    }
  }

  async model() {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients');
    const licenseStart = await this.getLicenseStartTime();
    const activity = await this.getActivity(licenseStart);
    return {
      config,
      versionHistory,
      activity,
      licenseStartTimestamp: licenseStart,
      currentDate: this.currentDate,
    };
  }
}
