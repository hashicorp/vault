/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';

import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  wait?: number; // defaults to 500
  onInput(value: string): void; // invoked with input value after debounce timer expires
}

export default class FilterInputComponent extends Component<Args> {
  @action
  onInput(event: HTMLElementEvent<HTMLInputElement>) {
    const wait = this.args.wait || 500;
    // ts complains when trying to pass object of optional args to callback as 3rd arg to debounce
    // eslint-disable-next-line
    // @ts-ignore
    debounce(this, this.args.onInput, event.target.value, wait);
  }
}
