/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithFormFieldsModel } from 'vault/app-types';

type PkiConfigCrlModel = WithFormFieldsModel & {
  autoRebuild: boolean;
  autoRebuildGracePeriod: string;
  enableDelta: boolean;
  expiry: string;
  deltaRebuildInterval: string;
  disable: boolean;
  ocspExpiry: string;
  ocspDisable: boolean;
  crlPath: string;
  get canSet(): boolean;
};

export default PkiConfigCrlModel;
