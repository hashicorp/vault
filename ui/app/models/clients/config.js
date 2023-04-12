/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  retentionMonths: [
    {
      validator: (model) => parseInt(model.retentionMonths) >= (model.reportingEnabled ? 24 : 0),
      message: (model) => {
        if (model.reportingEnabled) {
          return 'Retention period must be a minimum of 24 months.';
        }
        return 'Retention period must be greater than or equal to 0.';
      },
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

  @attr('string') enabled;

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
