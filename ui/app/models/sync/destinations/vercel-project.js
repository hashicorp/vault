/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationsBaseModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsVercelProjectModel extends SyncDestinationsBaseModel {
  @attr('string') accessToken;
  @attr('string') projectId;
  @attr('string') teamId;
  @attr('array') deploymentEnvironments;
}
