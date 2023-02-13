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
  @tracked warningDismissed;
  @tracked expiredDismissed;

  constructor() {
    super(...arguments);
    // If nothing is saved in localStorage or the user has updated their Vault version, do not dismiss any of the banners.
    const localStorageLicenseBannerObject = localStorage.getItem('licenseBanner');
    if (!localStorageLicenseBannerObject || localStorageLicenseBannerObject.version !== this.currentVersion) {
      localStorage.setItem('licenseBanner', { dismissType: '', version: this.currentVersion });
      return;
    }
    // if dismissType has previously been saved in localStorage, update tracked properties.
    this.setDismissType(localStorageLicenseBannerObject.dismissType);
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
  dismissBanner(dismissAction) {
    // dismissAction is either 'dismiss-warning' or 'dismiss-expired'
    const updatedLocalStorageObject = { dismissType: dismissAction, version: this.currentVersion };
    localStorage.setItem('licenseBanner', updatedLocalStorageObject);
    this.setDismissType(dismissAction);
  }

  setDismissType(dismissType) {
    // reset tracked properties to false
    this.warningDismissed = this.expiredDismissed = false;
    if (dismissType === 'dismiss-warning') {
      this.warningDismissed = true;
    } else if (dismissType === 'dismiss-expired') {
      this.expiredDismissed = true;
    } else {
      // if dismissType is empty do nothing.
      return;
    }
  }
}
