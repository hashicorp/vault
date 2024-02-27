/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

interface Args {
  attr?: AttrData;
  label?: string;
}
interface AttrData {
  name: string; // required if @attr is passed
  options?: {
    label?: string;
    helpText?: string;
    subText?: string;
    possibleValues?: string[];
  };
}

/**
 * @module EnableInput
 * EnableInput components render a disabled input with a hardcoded masked value beside an "Edit" button to "enable" the input.
 * Clicking "Edit" hides the disabled input and renders the yielded component. This way any data management is handled by the parent.
 * These are useful for editing inputs of sensitive values not returned by the API. The extra click ensures the user is intentionally editing the field.
 *
 * @example
 <EnableInput class="field" @attr={{attr}}>
  <FormField @attr={{attr}} @model={{@destination}} @modelValidations={{this.modelValidations}} />
 </EnableInput>

// without passing @attr
 <EnableInput @label="AWS password">
  <Input @type="text" />
 </EnableInput>

 * @param {object} [attr] - used to generate label for `ReadonlyFormField`, `name` key is required. Can be an attribute from a model exported with expandAttributeMeta.
 * @param {string} [label] - required if no attr passed. Used to ensure a11y conformance for the readonly input.
 */

export default class EnableInputComponent extends Component<Args> {
  @tracked enable = false;
}
