/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

export function secretQueryParams([backendType, type = ''], { asQueryParams }) {
  // Use effective engine type to handle external plugin mappings
  const effectiveBackendType = getEffectiveEngineType(backendType);
  const values = {
    transit: { tab: 'actions' },
    database: { type },
    keymgmt: { itemType: type === 'provider' ? 'provider' : 'key' },
  }[effectiveBackendType];
  // format required when using LinkTo with positional params
  if (values && asQueryParams) {
    return {
      isQueryParams: true,
      values,
    };
  }
  return values;
}

export default helper(secretQueryParams);
