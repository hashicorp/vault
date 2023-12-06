/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default class DiffController extends Controller.extend(BackendCrumbMixin) {}
