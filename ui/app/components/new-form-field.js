import Component from '@glimmer/component';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';

export default class NewFormFieldComponent extends Component {
  get value() {
    return this.args.model[this.args.name];
  }

  get label() {
    return this.args.attr.options?.label || capitalize([humanize([dasherize([this.args.attr.name])])]);
  }

  get editType() {
    if (this.args.attr.options.possibleValues) {
      return 'select';
    }
    console.log(this.args.attr.options);
    return this.args.attr.options?.editType || 'text';
  }
}
