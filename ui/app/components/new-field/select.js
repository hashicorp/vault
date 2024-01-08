import { assert } from '@ember/debug';
import { action } from '@ember/object';
import Component from '@glimmer/component';

export default class NewFieldSelectComponent extends Component {
  constructor() {
    super(...arguments);
    assert('new-field/select is missing required fields', this.args.name && this.args.label);
  }

  @action
  handleChange({ target }) {
    const { name, value } = target;
    this.args.onChange(name, value);
  }
}
