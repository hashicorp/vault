/**
 * Copyright IBM Corp. 2016, 2025
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
import type VersionService from 'vault/services/version';

enum Banners {
  EXPIRED = 'expired',
  WARNING = 'warning',
  PKI = 'pki-only-info',
}

interface Args {
  expiry: string;
  autoloaded: boolean;
}

/**
 * @module LicenseBanners
 * LicenseBanners components are used to display Vault-specific license messages
 *
 * @example
 * ```js
 * <LicenseBanners @expiry={expiryDate} />
 * ```
 * @param {string} expiry - RFC3339 date timestamp
 */

export default class LicenseBanners extends Component<Args> {
  @service declare readonly version: VersionService;

  @tracked warningDismissed = false;
  @tracked expiredDismissed = false;
  @tracked infoDismissed = false;

  banners = Banners;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    // reset and show a previously dismissed license banner if:
    // the version has been updated or the license has been updated (indicated by a change in the expiry date).
    const item = localStorage.getItem(this.dismissedBannerKey) ?? []; // returns warning, expired and/or pki-only-info
    // older entries will not be an array as it was either "expired" OR "warning"
    // with the addition of "pki-only-info", it can hold all values
    // this check maintains backwards compatibility with the previous format
    const bannerTypes = Array.isArray(item) ? item : [item];

    bannerTypes.forEach((type) => {
      this.updateDismissType(type);
    });
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

  get isPKIOnly() {
    return this.version.hasPKIOnly;
  }

  @action
  dismissBanner(dismissAction: Banners) {
    // if a client's version changed their old localStorage key will still exists.
    localStorage.cleanupStorage('dismiss-license-banner', this.dismissedBannerKey);

    // updates localStorage and then updates the template by calling updateDismissType
    const item = localStorage.getItem(this.dismissedBannerKey) ?? [];
    // older entries will not be an array as it was either "expired" OR "warning"
    // with the addition of "pki-only-info", it can hold all values
    // this check maintains backwards compatibility with the previous format
    const bannerTypes = Array.isArray(item) ? item : [item];
    localStorage.setItem(this.dismissedBannerKey, [...bannerTypes, dismissAction]);

    this.updateDismissType(dismissAction);
  }

  updateDismissType(dismissType?: Banners) {
    // updates tracked properties to update template
    switch (dismissType) {
      case this.banners.WARNING:
        this.warningDismissed = true;
        break;
      case this.banners.EXPIRED:
        this.expiredDismissed = true;
        break;
      case this.banners.PKI:
        this.infoDismissed = true;
        break;
      default:
        break;
    }
  }
}
