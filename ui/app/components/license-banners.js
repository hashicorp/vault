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
  @tracked dismissLicenseExpired = false;
  @tracked localStorageLicenseBannerState = localStorage.getItem('licenseBannerState');

  constructor() {
    super(...arguments);
    if (!this.localStorageLicenseBannerState) {
      localStorage.setItem('licenseBannerState', { dismiss: false, version: this.currentVersion });
    } else {
      if (
        this.localStorageLicenseBannerState.version === this.currentVersion &&
        this.localStorageLicenseBannerState.dismiss
      ) {
        this.dismissLicenseExpired = true;
      }
    }
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
  dismissLicenseExpiredBanner() {
    const updatedObject = { dismiss: true, version: this.currentVersion };
    localStorage.setItem('licenseBannerState', updatedObject);
    this.dismissLicenseExpired = true;
  }
}
