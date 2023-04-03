/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import isAfter from 'date-fns/isAfter';
import differenceInDays from 'date-fns/differenceInDays';
import localStorage from 'vault/lib/local-storage';
import timestamp from 'core/utils/timestamp';

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

  @tracked warningDismissed;
  @tracked expiredDismissed;

  constructor() {
    super(...arguments);
    // do not dismiss any banners if the user has updated their version
    const dismissedBanner = localStorage.getItem(`dismiss-license-banner-${this.currentVersion}`); // returns either warning or expired
    this.updateDismissType(dismissedBanner);
  }

  get currentVersion() {
    return this.version.version;
  }

  get licenseExpired() {
    if (!this.args.expiry) return false;
    return isAfter(timestamp.now(), new Date(this.args.expiry));
  }

  get licenseExpiringInDays() {
    // Anything more than 30 does not render a warning
    if (!this.args.expiry) return 99;
    return differenceInDays(new Date(this.args.expiry), timestamp.now());
  }

  @action
  dismissBanner(dismissAction) {
    // if a client's version changed their old localStorage key will still exists.
    localStorage.cleanUpStorage('dismiss-license-banner', `dismiss-license-banner-${this.currentVersion}`);
    // updates localStorage and then updates the template by calling updateDismissType
    localStorage.setItem(`dismiss-license-banner-${this.currentVersion}`, dismissAction);
    this.updateDismissType(dismissAction);
  }

  updateDismissType(dismissType) {
    // updates tracked properties to update template
    if (dismissType === 'warning') {
      this.warningDismissed = true;
    } else if (dismissType === 'expired') {
      this.expiredDismissed = true;
    }
  }
}
