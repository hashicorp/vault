/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { allFeatures } from 'vault/helpers/all-features';
/**
 * @module LicenseInfo
 *
 * @example
 * ```js
 * <LicenseInfo
 *   @startTime="2020-03-12T23:20:50.52Z"
 *   @expirationTime="2021-05-12T23:20:50.52Z"
 *   @licenseId="some-license-id"
 *   @features={{array 'Namespaces' 'DR Replication'}}
 *   @autoloaded={{true}}
 *   @performanceStandbyCount=1
 * />
 *
 * @param {string} startTime - RFC3339 formatted timestamp of when the license became active
 * @param {string} expirationTime - RFC3339 formatted timestamp of when the license will expire
 * @param {string} licenseId - unique ID of the license
 * @param {Array<string>} features - Array of feature names active on license
 * @param {boolean} autoloaded - Whether the license is autoloaded
 * @param {number} performanceStandbyCount - Number of performance standbys active
 */
export default class LicenseInfoComponent extends Component {
  get featuresInfo() {
    return allFeatures().map((feature) => {
      const active = this.args.features.includes(feature);
      if (active && feature === 'Performance Standby') {
        const count = this.args.performanceStandbyCount;
        return {
          name: feature,
          active: count ? active : false,
          count,
        };
      }
      return { name: feature, active };
    });
  }
}
