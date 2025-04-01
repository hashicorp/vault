/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

interface AfterOptionsArgs {
  loadOptions: () => void;
}

/**
 * @module AfterOptions
 * @description component is used to display action items inside the namespace picker dropdown.
 *  The "Manage" button directs the user to the namespace management page.
 *  The "Refresh List" button refrehes the list of namespaces in the dropdown.
 *
 * @example
 * @afterOptionsComponent={{component "namespace-picker/after-options" loadOptions=this.loadOptions}}
 */

export default class AfterOptions extends Component<AfterOptionsArgs> {
  @action
  refreshNamespaceList() {
    this.args?.loadOptions();
  }
}
