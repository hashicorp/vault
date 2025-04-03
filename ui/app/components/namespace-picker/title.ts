/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

interface TitleArgs {
  all: boolean;
  count: () => void;
}

/**
 * @module Title
 * @description Component used to display title in the namespace-picker.
 * By default, display "All Namespaces" with the total count.
 * When a search is entered, display "Matching Namespaces" with the count of matching results.
 *
 * @example
 * {{component "namespace-picker/title" count=this.options.count}}
 */

export default class Title extends Component<TitleArgs> {
  @tracked count = this.args?.count;
  @tracked scope = this.args?.all ? 'All' : 'Matching';
}
