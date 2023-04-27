/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action, set } from '@ember/object';

/**
 * @module TransformAdvancedTemplating
 * TransformAdvancedTemplating components are used to modify encode/decode formats of transform templates
 *
 * @example
 * ```js
 * <TransformAdvancedTemplating @model={{this.model}} />
 * ```
 * @param {Object} model - transform template model
 */

export default class TransformAdvancedTemplating extends Component {
  @tracked inputOptions = [];

  @action
  setInputOptions(testValue, captureGroups) {
    if (captureGroups && captureGroups.length) {
      this.inputOptions = captureGroups.map(({ position, value }) => {
        return {
          label: `${position}: ${value}`,
          value: position,
        };
      });
    } else {
      this.inputOptions = [];
    }
  }
  @action
  decodeFormatValueChange(kvObject, kvData, value) {
    set(kvObject, 'value', value);
    this.args.model.decodeFormats = kvData.toJSON();
  }
}
