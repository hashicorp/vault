/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module CertificateCard
 * The CertificateCard component receives data and optionally receives a boolean declaring if that data is meant to be in PEM
 * Format. It renders using the <HDS::Card::Container>. To the left there is a certificate icon. In the center there is a label
 * which says which format (PEM or DER) the data is in. Below the label is the truncated data. To the right there is a copy
 * button to copy the data.
 *
 * @example
 * ```js
 *  <CertificateCard @data={{value}} @isPem={{true}} />
 * ```
 * @param {string} data - the data to be displayed in the component (usually in PEM or DER format)
 * @param {boolean} isPem - optional argument for if the data is required to be in PEM format (and should thus have the PEM Format label)
 */

export default class CertificateCardComponent extends Component {
  // Returns the format the data is in: PEM, DER, or no format if no data is provided
  get format() {
    if (!this.args.data) return '';

    let value;
    if (typeof this.args.data === 'object') {
      value = this.args.data[0];
    } else {
      value = this.args.data;
    }

    if (value.substring(0, 11) === '-----BEGIN ' || this.args.isPem === true) {
      return 'PEM Format';
    }
    return 'DER Format';
  }
}
