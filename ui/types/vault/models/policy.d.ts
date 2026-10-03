/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

export default class PolicyModel extends Model {
  id: string;
  name: string;
  policy: string;
  policyType: string;
  format: 'json' | 'hcl';
}
