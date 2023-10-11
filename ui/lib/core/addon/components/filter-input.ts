/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce, next } from '@ember/runloop';

import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  value?: string; // initial value
  placeholder?: string; // defaults to Type to filter results
  wait?: number; // defaults to 200
  autofocus?: boolean; // initially focus the input on did-insert
  onInput(value: string): void;
}

export default class FilterInputComponent extends Component<Args> {
  value: string | undefined;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.value = this.args.value;
  }

  get placeholder() {
    return this.args.placeholder || 'Type to filter results';
  }

  @action
  focus(elem: HTMLElement) {
    if (this.args.autofocus) {
      next(() => elem.focus());
    }
  }

  @action
  onInput(event: HTMLElementEvent<HTMLInputElement>) {
    const callback = () => {
      this.args.onInput(event.target.value);
    };
    const wait = this.args.wait || 200;
    // ts complains when trying to pass object of optional args to callback as 3rd arg to debounce
    debounce(this, callback, wait);
  }
}
