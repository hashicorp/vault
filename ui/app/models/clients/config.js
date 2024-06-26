/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  retentionMonths: [
    {
      validator: (model) => parseInt(model.retentionMonths) >= model.minimumRetentionMonths,
      message: (model) =>
        `Retention period must be greater than or equal to ${model.minimumRetentionMonths}.`,
    },
    {
      validator: (model) => parseInt(model.retentionMonths) <= 60,
      message: 'Retention period must be less than or equal to 60.',
    },
  ],
};

@withModelValidations(validations)
@withFormFields(['enabled', 'retentionMonths'])
export default class ClientsConfigModel extends Model {
  @attr('boolean') queriesAvailable; // true only if historical data exists, will be false if there is only current month data

  @attr('number', {
    label: 'Retention period',
    subText: 'The number of months of activity logs to maintain for client tracking.',
  })
  retentionMonths;

  @attr('number') minimumRetentionMonths;

  // refers specifically to the activitylog and will always be on for enterprise
  @attr('string') enabled;

  // reporting_enabled is for automated reporting and only true of the customer hasn’t opted-out of automated license reporting
  @attr('boolean') reportingEnabled;

  @attr('date') billingStartTimestamp;

  @lazyCapabilities(apiPath`sys/internal/counters/config`) configPath;

  get canRead() {
    return this.configPath.get('canRead') !== false;
  }
  get canEdit() {
    return this.configPath.get('canUpdate') !== false;
  }
}
