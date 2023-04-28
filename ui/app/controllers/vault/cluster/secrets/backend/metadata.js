/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import { action } from '@ember/object';

export default class MetadataController extends Controller.extend(BackendCrumbMixin) {
  @action
  refreshModel() {
    this.send('refreshModel');
  }
}
