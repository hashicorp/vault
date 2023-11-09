/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module ReadonlyFormField
 * ReadonlyFormField components are used to display read only, non-editable attributes
 *
 * @example
 * <ReadonlyFormField @attr={{hash name="my attr"}} @value="some value" />
 *
 * @param {object} attr - Should be an attribute from a model exported with expandAttributeMeta
 * @param {any} value - The value that should be displayed on the input
 */

import Component from '@glimmer/component';
import { setComponentTemplate } from '@ember/component';
import { capitalize, dasherize } from '@ember/string';
import { humanize } from 'vault/helpers/humanize';
import layout from '../templates/components/readonly-form-field';

class ReadonlyFormField extends Component {
  get labelString() {
    if (!this.args.attr) {
      return '';
    }
    const label = this.args.attr.options ? this.args.attr.options.label : '';
    const name = this.args.attr.name;
    if (label) {
      return label;
    }
    if (name) {
      return capitalize(humanize([dasherize(name)]));
    }
    return '';
  }
}

export default setComponentTemplate(layout, ReadonlyFormField);
