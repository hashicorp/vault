/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import type NamespaceService from 'vault/services/namespace';
import type VersionService from 'vault/services/version';

export default class DashboardVaultVersionTitle extends Component {
  @service declare readonly namespace: NamespaceService;
  @service declare readonly version: VersionService;
}
