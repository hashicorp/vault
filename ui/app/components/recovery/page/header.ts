/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import type VersionService from 'vault/services/version';

interface Args {
  title: string;
  subtitle?: string;
  action?: unknown;
}

export default class Header extends Component<Args> {
  @service declare readonly version: VersionService;
}
