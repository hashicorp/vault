/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module CertificateCard
 * [Description]
 *
 * @example
 * ```js
 *  <CertificateCard @certificateValue={{value}} @isPem={{true}} />
 * ```
 * @param {string} certificateValue - the value to be displayed
 * @param {boolean} isPem - optional argument for if the certificateValue is required to be PEM format
 */

export default class CertificateCardComponent extends Component {
  get format() {
    if (!this.args.certificateValue) return '';

    let value;
    if (typeof this.args.certificateValue === 'object') {
      value = this.args.certificateValue[0];
    } else {
      value = this.args.certificateValue;
    }

    if (value.substring(0, 11) === '-----BEGIN ' || this.args.isPem === true) {
      return 'PEM Format';
    }
    return 'DER Format';
  }
}
