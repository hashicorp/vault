/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

interface Args {
  attr: Attr;
}
interface Attr {
  name: string;
  options: {
    label: string;
    helpText: string;
    subText: string;
    possibleValues: string[];
  };
}

/**
 * @module EnableInput
 * EnableInput components wrap a form field component and include an "Edit" button that users must click to enable the input.
 * These are useful for inputs sensitive values that are not returned by the API, so they are only sent in a POST request if the user
 * has performed an extra click to intentionally edit the field.
 *
 * @example
 * <EnableInput class="field" @attr={{attr}}>
 *  <FormField @attr={{attr}} @model={{@destination}} @modelValidations={{this.modelValidations}} />
 * </EnableInput>
 *
 */

export default class EnableInputComponent extends Component<Args> {
  @tracked enable = false;
}
