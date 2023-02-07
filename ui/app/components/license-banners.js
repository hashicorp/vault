/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

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

import Component from '@glimmer/component';
import isAfter from 'date-fns/isAfter';
import differenceInDays from 'date-fns/differenceInDays';

export default class LicenseBanners extends Component {
  get licenseExpired() {
    if (!this.args.expiry) return false;
    return isAfter(new Date(), new Date(this.args.expiry));
  }

  get licenseExpiringInDays() {
    // Anything more than 30 does not render a warning
    if (!this.args.expiry) return 99;
    return differenceInDays(new Date(this.args.expiry), new Date());
  }
}
