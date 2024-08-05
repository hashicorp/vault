/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { FormField } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';

export default interface AwsRootConfigModel extends WithFormFieldsAndValidationsModel {
  backend: string;
  name: string;
  secret_key: string;
  access_key: string;
  region: string;
  iam_endpoint: string;
  sts_endpoint: string;
  max_retries: number;
}
