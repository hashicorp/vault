/**
 * @module ReadonlyFormField
 * ReadonlyFormField components are used to...
 *
 * @example
 * ```js
 * <ReadonlyFormField @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} attr - Should be an attribute from a model exported with expandAttributeMeta
 * @param {any} value - This value's type should differ depending on
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import { setComponentTemplate } from '@ember/component';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';
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
      return capitalize([humanize([dasherize([name])])]);
    }
    return '';
  }
}

export default setComponentTemplate(layout, ReadonlyFormField);
