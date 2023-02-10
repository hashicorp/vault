import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import isAfter from 'date-fns/isAfter';
import differenceInDays from 'date-fns/differenceInDays';
import localStorage from 'vault/lib/local-storage';

/**
 * @module LicenseBanners
 * LicenseBanners components are used to display Vault-specific license expiry messages
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
  @tracked localStorageLicenseBannerObject = localStorage.getItem('licenseBanner');
  @tracked dismissType = 'none'; // dismissType options: both, dismiss-warning, dismiss-expired, none

  constructor() {
    super(...arguments);
    // if nothing is saved in localStorage, or the user has updated their Vault version show the license banners
    if (
      !this.localStorageLicenseBannerObject ||
      this.localStorageLicenseBannerObject.version !== this.currentVersion
    ) {
      localStorage.setItem('licenseBanner', { dismissType: 'none', version: this.currentVersion });
    }

    // update tracked property to equal either dismissType from localStorage or 'none' if the local storage object does not exists.
    this.dismissType = !this.localStorageLicenseBannerObject?.dismissType
      ? 'none'
      : this.localStorageLicenseBannerObject.dismissType;
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

  get showWarning() {
    return this.dismissType === 'dismiss-expired' || this.dismissType === 'none' ? true : false;
  }

  get showExpired() {
    return this.dismissType === 'dismiss-warning' || this.dismissType === 'none' ? true : false;
  }

  @action
  dismissBanner(dismissAction) {
    // dismissAction is either 'dismiss-warning' or 'dismiss-expired'
    const updatedLocalStorageObject =
      this.dismissType === 'none'
        ? { dismissType: dismissAction, version: this.currentVersion }
        : { dismissType: 'both', version: this.currentVersion };

    localStorage.setItem('licenseBanner', updatedLocalStorageObject);
    // update tracked property so showWarning and showExpired are updated.
    this.dismissType = this.dismissType === 'none' ? dismissAction : 'both';
  }
}
