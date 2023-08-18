/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class InputSearch extends Component {
  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   *
   */
  @action
  inputChanged(e) {
    this.args.onChange(e.target.value);
  }
}
