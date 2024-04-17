/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import type PkiIssuerModel from 'vault/models/pki/issuer';

interface Args {
  issuer: PkiIssuerModel;
}

export default class PkiIssuerDetailsComponent extends Component<Args> {
  @tracked showRotationModal = false;

  get parsingErrors() {
    if (this.args.issuer.parsedCertificate?.parsing_errors?.length) {
      return this.args.issuer.parsedCertificate.parsing_errors.map((e: Error) => e.message).join(', ');
    }
    return '';
  }
}
