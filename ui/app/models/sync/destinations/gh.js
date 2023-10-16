/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationsBaseModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsGithubModel extends SyncDestinationsBaseModel {
  @attr('string') accessToken;
  @attr('string') repositoryOwner;
  @attr('string') repositoryName;
}
