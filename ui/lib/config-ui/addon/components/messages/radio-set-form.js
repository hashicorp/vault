/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display list of messages.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @messages={{this.messages}}  />
 * ```
 * @param {array} messages - array message objects
 */

export default class RadioSetForm extends Component {
  @action
  updateRadioValue(evt) {
    this.args.model[evt.target.name] = evt.target.value;
  }
}
