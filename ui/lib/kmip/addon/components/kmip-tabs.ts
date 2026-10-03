/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import type RouterService from '@ember/routing/router-service';

export default class KmipTabs extends Component {
  @service('app-router') declare readonly router: RouterService;
}
