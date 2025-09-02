/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';
import type PkiKeyModel from 'vault/models/pki/key';

interface Args {
  keyModels: PkiKeyModel[];
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
