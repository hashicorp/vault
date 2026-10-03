/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

interface Args {
  keys: { key_id: string; is_default: boolean; key_name: string }[];
  mountPoint: string;
  backend: string;
  canImportKey: boolean;
  canGenerateKey: boolean;
  canRead: boolean;
  canEdit: boolean;
  hasConfig: boolean;
}

export default class PkiKeyList extends Component<Args> {
  notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
}
