/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module InputSearch
 * This component renders an input that fires a callback on "keyup" containing the input's value
 *
 * @example
 * <InputSearch
 * @initialValue="secret/path/"
 * @onChange={{this.handleSearch}}
 * @placeholder="search..."
 * />
 * @param {string} [id] - unique id for the input
 * @param {string} [initialValue] - initial search value, i.e. a secret path prefix, that pre-fills the input field
 * @param {string} [placeholder] - placeholder text for the input
 * @param {string} [label] - label for the input
 * @param {string} [subtext] - displays below the label
 */

export default class InputSearch extends Component {
  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   *
   */
  @tracked searchInput = '';

  constructor() {
    super(...arguments);
    this.searchInput = this.args?.initialValue;
  }

  @action
  inputChanged() {
    this.args.onChange(this.searchInput);
  }
}
