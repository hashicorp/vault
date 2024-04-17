/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { debounce, next } from '@ember/runloop';

import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  wait?: number; // defaults to 500
  autofocus?: boolean; // initially focus the input on did-insert
  hideIcon?: boolean; // hide the search icon in the input
  onInput(value: string): void; // invoked with input value after debounce timer expires
}

export default class FilterInputComponent extends Component<Args> {
  @action
  focus(elem: HTMLElement) {
    if (this.args.autofocus) {
      next(() => elem.focus());
    }
  }

  @action
  onInput(event: HTMLElementEvent<HTMLInputElement>) {
    const wait = this.args.wait || 500;
    // ts complains when trying to pass object of optional args to callback as 3rd arg to debounce
    // eslint-disable-next-line
    // @ts-ignore
    debounce(this, this.args.onInput, event.target.value, wait);
  }
}
