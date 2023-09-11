/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';

import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  placeholder?: string; // defaults to Type to filter results
  wait?: number; // defaults to 200
  onInput(value: string): void;
}

export default class FilterInputComponent extends Component<Args> {
  get placeholder() {
    return this.args.placeholder || 'Type to filter results';
  }

  @action onInput(event: HTMLElementEvent<HTMLInputElement>) {
    const callback = () => {
      this.args.onInput(event.target.value);
    };
    const wait = this.args.wait || 200;
    // ts complains when trying to pass object of optional args to callback as 3rd arg to debounce
    debounce(this, callback, wait);
  }
}
