/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type VersionService from 'vault/services/version';

export default class SyncLandingCtaComponent extends Component {
  @service declare readonly version: VersionService;
}
