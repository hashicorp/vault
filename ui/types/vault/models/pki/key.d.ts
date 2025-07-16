/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithValidationsModel } from 'vault/app-types';

type PkiKeyModel = WithValidationsModel & {
  secretMountPath: class;
  keyId: string;
  keyName: string;
  type: string;
  keyType: string;
  keyBits: string;
  pemBundle: string;
  privateKey: string;
  get backend(): string;
  get canRead(): boolean;
  get canEdit(): boolean;
  get canDelete(): boolean;
  get canGenerateKey(): boolean;
  get canImportKey(): boolean;
};

export default PkiKeyModel;
