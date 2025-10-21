/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { WithFormFieldsAndValidationsModel } from 'vault/vault/app-types';

export default interface ClientsConfigModel extends WithFormFieldsAndValidationsModel {
  queriesAvailable: boolean;
  retentionMonths: number;
  minimumRetentionMonths: number;
  enabled: string;
  reportingEnabled: boolean;
  billingStartTimestamp: Date;
}
