/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

// import type { ModelValidations, FormField, FormFieldGroups } from 'vault/app-types';
// import type MountConfigModel from 'vault/models/mount-config';

export default class LicenseModel extends Model {
  features: Array<string>;
}
