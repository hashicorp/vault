/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationsBaseModel from './base';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsVercelProjectModel extends SyncDestinationsBaseModel {
  @attr('string') accessToken;
  @attr('string') projectId;
  @attr('string') teamId;
  @attr('array') deploymentEnvironments;

  get type() {
    return 'vercel-project';
  }

  get icon() {
    return 'vercel';
  }
}
