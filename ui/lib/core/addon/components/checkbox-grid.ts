/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'forms';

interface CheckboxGridArgs {
  name: string;
  label: string;
  subText?: string;
  fields: Field[];
  value: string[] | undefined;
  onChange: (name: string, value: string[]) => void;
}
interface Field {
  key: string;
  label: string;
}

/**
 * @module CheckboxGrid
 * CheckboxGrid components are used to allow users to select any
 * number of predetermined options, aligned in a 3-column grid.
 *
 * @example
 * ```js
 * <CheckboxGrid
 *   @name="modelKey"
 *   @label="Model Attribute Label"
 *   @fields={{options}}
 *   @value={{['Hello', 'Yes']}}
 * />
 * ```
 */

export default class CheckboxGrid extends Component<CheckboxGridArgs> {
  get checkboxes() {
    const list = this.args.value || [];
    return this.args.fields.map((field) => ({
      ...field,
      value: list.includes(field.key),
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
