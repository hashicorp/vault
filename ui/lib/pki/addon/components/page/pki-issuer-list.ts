/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

import type { PkiReadIssuerResponse } from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';

type Issuer = PkiReadIssuerResponse & {
  id: string;
  is_default: boolean;
  serial_number: string;
  isRoot: boolean;
  parsedCertificate: ParsedCertificateData;
};
interface Args {
  issuers: Issuer[];
  mountPoint: string;
  backend: string;
}

export default class PkiIssuerList extends Component<Args> {
  notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
}
