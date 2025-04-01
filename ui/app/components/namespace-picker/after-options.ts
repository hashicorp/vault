/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

interface AfterOptionsArgs {
  loadOptions: () => void;
}

export default class AfterOptions extends Component<AfterOptionsArgs> {
  @action
  refreshNamespaceList() {
    this.args?.loadOptions();
  }
}
