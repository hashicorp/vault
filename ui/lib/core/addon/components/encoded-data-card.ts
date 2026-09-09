/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module EncodedDataCard
 * The EncodedDataCard component renders truncated encoded data with a copy button.
 * It detects if the format is PEM or optionally receives a boolean explicitly stating PEM.
 * If @data is not already a string, it converts it into a copyable string.
 *
 * @example
 *  <EncodedDataCard @isPem={{true}} @data="-----BEGIN CERTIFICATE-----MIIDezCCAmOgAwIBAgIUTBbQcZijQsmd0rjd6COikPsrGyowDQYJKoZIhvcNAQELBQAwFDESMBAGA1UEAxMJdGVzdC1yb290MB4XDTIzMDEyMDE3NTcxMloXDTIzMDIy" />
 *
 * @param {string|Array|object} data - the encoded data to be displayed in the component, such as PEM, DER, or a base64-encoded string
 * @param {boolean} [isPem] - optional argument for if the data should be labeled as PEM format
 */

interface Args {
  data?: string | string[] | Record<string, unknown> | unknown[];
  isPem?: boolean;
}

export default class EncodedDataCardComponent extends Component<Args> {
  get certDisplay(): { label: string; icon: string } {
    if (!this.args.data) return { label: '', icon: 'transform-data' };

    const value = Array.isArray(this.args.data) ? this.args.data[0] : this.args.data;
    const hasPemHeader = typeof value === 'string' && value.substring(0, 11) === '-----BEGIN ';
    if (hasPemHeader || this.args.isPem === true) {
      return { label: 'PEM Format', icon: 'certificate' };
    }

    return { label: 'Encoded Data', icon: 'transform-data' };
  }

  get copyValue(): string {
    const { data } = this.args;
    if (!data) return '';

    if (Array.isArray(data)) {
      return data
        .map((value) => {
          if (typeof value === 'string') return value;
          if (value && typeof value === 'object') return JSON.stringify(value);
          return String(value);
        })
        .join('\n');
    }

    switch (typeof data) {
      case 'string':
        return data;
      case 'object':
        // unlikely for certificates but just in case
        return JSON.stringify(data);
      default:
        return String(data);
    }
  }
}
