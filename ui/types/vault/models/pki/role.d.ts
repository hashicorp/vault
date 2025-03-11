/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithValidationsModel } from 'vault/app-types';

type PkiRoleModel = WithValidationsModel & {
  name: string;
  issuerRef: string;
  keyType: string;
  keyBits: string | undefined;
};

export default PkiRoleModel;
