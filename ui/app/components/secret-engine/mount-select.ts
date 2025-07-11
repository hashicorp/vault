/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';

/**
 * @module MountSelect todo
 */

interface Args {
  mountsArray: Array<string>;
}

export default class SecretEngineMountSelect extends Component<Args> {
  @service declare readonly router: RouterService;

  @action selectMount(type: string) {
    this.router.transitionTo('vault.cluster.secrets.mount.create', type);
  }
}
