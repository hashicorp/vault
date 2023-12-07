/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

interface Args {
  attrOptions?: AttrOptions;
}
interface AttrOptions {
  name: string; // required
  options?: {
    label?: string;
    helpText?: string;
    subText?: string;
    possibleValues?: string[];
  };
}

/**
 * @module EnableInput
 * EnableInput components render a disabled input with a hardcoded value of ********** with an "Edit" button to "enable" the input,
 * clicking "Edit" hides the disabled input and instead renders the yielded component. This way any data management is handled by the parent.
 * These are useful for inputs sensitive values that are not returned by the API, so they are only sent in a POST request if the user
 * has performed an extra click to intentionally edit the field.
 *
 * @example
 * <EnableInput class="field" @attr={{attr}}>
 *  <FormField @attr={{attr}} @model={{@destination}} @modelValidations={{this.modelValidations}} />
 * </EnableInput>
 *
 * <EnableInput class="field" @attr={{attr}}>
 *  <FormField @attr={{attr}} @model={{@destination}} @modelValidations={{this.modelValidations}} />
 * </EnableInput>
 *
 * @param {object} [attrOptions] - used to generate the label for `ReadonlyFormField`, `name` key is required. Can be an attribute from a model exported with expandAttributeMeta.
 */

export default class EnableInputComponent extends Component<Args> {
  @tracked enable = false;
}
