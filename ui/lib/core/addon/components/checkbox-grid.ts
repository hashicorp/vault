/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
 * @description
 * CheckboxGrid components are used to allow users to select any number of predetermined options, aligned in a 3-column grid.
 *
 *
 * @example
 * <CheckboxGrid @name="extKeyUsage" @label="Extended key usage" @fields={{array (hash key="EmailProtection" label="Email Protection") (hash key="TimeStamping" label="Time Stamping") (hash key="ServerAuth" label="Server Auth") }} @value={{array "TimeStamping"}} />
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
