import { assert } from '@ember/debug';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { HTMLElementEvent } from 'vault/forms';

interface Args {
  name: string;
  label: string;
  onChange: CallableFunction;
  value: string;
  fieldErrors?: string[];
  isRequired?: boolean;
  type?: string;
}

export default class NewFieldTextInputComponent extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(
      'new-field/text-input is missing required fields',
      this.args.name && this.args.label && this.args.onChange
    );
  }

  get isValid() {
    return !this.args.fieldErrors || this.args.fieldErrors.length === 0;
  }

  @action
  handleChange({ target }: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = target;
    this.args.onChange(name, value);
  }
}
