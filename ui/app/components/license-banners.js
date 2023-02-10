import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import isAfter from 'date-fns/isAfter';
import differenceInDays from 'date-fns/differenceInDays';
import localStorage from 'vault/lib/local-storage';

/**
 * @module LicenseBanners
 * LicenseBanners components are used to display Vault-specific license expiry messages.
 *
 * @example
 * ```js
 * <LicenseBanners @expiry={expiryDate} />
 * ```
 * @param {string} expiry - RFC3339 date timestamp
 */

export default class LicenseBanners extends Component {
  @service version;

  @tracked currentVersion = this.version.version;
  @tracked localStorageLicenseBannerState = localStorage.getItem('licenseBannerState');
  @tracked dismissType = 'none';

  // dismissType options = [both, warning, expired, none]
  constructor() {
    super(...arguments);
    if (
      !this.localStorageLicenseBannerState ||
      this.localStorageLicenseBannerState.version !== this.currentVersion
    ) {
      localStorage.setItem('licenseBannerState', { dismissType: 'none', version: this.currentVersion });
    }

    this.dismissType = !this.localStorageLicenseBannerState?.dismissType
      ? 'none'
      : this.localStorageLicenseBannerState.dismissType;
  }

  get showWarning() {
    return this.dismissType === 'both' || this.dismissType === 'dismiss-warning' ? false : true;
  }
  get showExpired() {
    return this.dismissType === 'both' || this.dismissType === 'dismiss-expired' ? false : true;
  }

  get licenseExpired() {
    if (!this.args.expiry) return false;
    return isAfter(new Date(), new Date(this.args.expiry));
  }

  get licenseExpiringInDays() {
    // Anything more than 30 does not render a warning
    if (!this.args.expiry) return 99;
    return differenceInDays(new Date(this.args.expiry), new Date());
  }

  @action
  dismissBanner(bannerType) {
    const updatedLicenseBannerState =
      this.dismissType === 'none'
        ? { dismissType: bannerType, version: this.currentVersion }
        : { dismissType: 'both', version: this.currentVersion };

    localStorage.setItem('licenseBannerState', updatedLicenseBannerState);

    this.dismissType = this.dismissType === 'none' ? bannerType : 'both';
  }
}
