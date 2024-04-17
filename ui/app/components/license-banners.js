/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
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
    // reset and show a previously dismissed license banner if:
    // the version has been updated or the license has been updated (indicated by a change in the expiry date).
    const bannerType = localStorage.getItem(this.dismissedBannerKey); // returns either warning or expired

    this.updateDismissType(bannerType);
  }

  get currentVersion() {
    return this.version.version;
  }

  get dismissedBannerKey() {
    return `dismiss-license-banner-${this.currentVersion}-${this.args.expiry}`;
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
    localStorage.cleanupStorage('dismiss-license-banner', this.dismissedBannerKey);
    // updates localStorage and then updates the template by calling updateDismissType
    localStorage.setItem(this.dismissedBannerKey, dismissAction);
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
