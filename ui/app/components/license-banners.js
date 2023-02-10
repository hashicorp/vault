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
  @tracked localStorageLicenseBannerState = localStorage.getItem('licenseBannerState');
  @tracked dismissType = 'none'; // dismissType options: [both, dismiss-warning, dismiss-expired, none]

  constructor() {
    super(...arguments);
    // if nothing saved in localStorage or the user has updated their version show both license banners
    if (
      !this.localStorageLicenseBannerState ||
      this.localStorageLicenseBannerState.version !== this.currentVersion
    ) {
      localStorage.setItem('licenseBannerState', { dismissType: 'none', version: this.currentVersion });
    }

    // set tracked property dismissType from either the local storage object or 'none' if one does not exist.
    this.dismissType = !this.localStorageLicenseBannerState?.dismissType
      ? 'none'
      : this.localStorageLicenseBannerState.dismissType;
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
  dismissBanner(bannerType) {
    // bannerType is either 'dismiss-warning' or 'dismiss-expired'
    const updatedLocalStorageObject =
      this.dismissType === 'none'
        ? { dismissType: bannerType, version: this.currentVersion }
        : { dismissType: 'both', version: this.currentVersion };

    localStorage.setItem('licenseBannerState', updatedLocalStorageObject);
    // update tracked property so showWarning and showExpired are updated causing the template to hide the appropriate banner.
    this.dismissType = this.dismissType === 'none' ? bannerType : 'both';
  }
}
