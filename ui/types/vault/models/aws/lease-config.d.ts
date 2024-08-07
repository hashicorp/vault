/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { FormField } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';

export default interface AwsRootLeaseModel extends WithFormFieldsAndValidationsModel {
  backend: string;
  leaseMax: string;
  lease: string;
}
