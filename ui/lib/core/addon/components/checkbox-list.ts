import Component from '@glimmer/component';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'forms';

interface Field {
  key: string;
  label: string;
}
interface CheckboxListArgs {
  name: string;
  label: string;
  subText?: string;
  fields: Field[];
  value: string[] | undefined;
  onChange: (name: string, value: string[]) => void;
}

/**
 * @module CheckboxList
 * CheckboxList components are used to allow users to select any number of predetermined options.
 *
 * @example
 * ```js
 * <CheckboxList @name="modelKey" @label="Model Attribute Label" @fields={{options}} @value={{['Hello', 'Yes']}}/>
 * ```
 */

export default class CheckboxList extends Component<CheckboxListArgs> {
  get checkboxes() {
    const list = this.args.value || [];
    return this.args.fields.map((field) => ({
      ...field,
      value: list.indexOf(field.key) >= 0 ? true : false,
    }));
  }

  @action checkboxChange(event: HTMLElementEvent<HTMLInputElement>) {
    const list = this.args.value || [];
    const checkboxName = event.target.id;
    const checkboxVal = event.target.checked;
    const idx = list.indexOf(checkboxName);
    if (checkboxVal === true && idx < 0) {
      list.push(checkboxName);
    } else if (checkboxVal === false && idx >= 0) {
      list.splice(idx, 1);
    }
    this.args.onChange(this.args.name, list);
  }
}
