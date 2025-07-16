/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import type { ValidationMap } from 'vault/vault/app-types';

/**
 * @module FormFieldGroups
 * FormFieldGroups components are field groups associated with a particular model. They render individual FormField components.
 *
 * @example
 * <FormFieldGroups @model={{mountModel}} @onChange={{action "onTypeChange"}} @renderGroup="Method Options" @onKeyUp={{action "onKeyUp"}} @validationMessages={{validationMessages}} />
 *
 * @param {Model} model- Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered.
 * @param {string} [renderGroup] - An allow list of groups to include in the render.
 * @param {function} [onChange] - Handler that will get set on the `FormField` component.
 * @param {function} [onKeyUp] - Handler that will set the value and trigger validation on input changes
 * @param {object} [modelValidations] - Object containing validation message for each property
 * @param {string} [groupName=fieldGroups] - attribute name where the field groups are
 * @param {boolean} [useEnableInput=false] - if you want to wrap your sensitive fields with the EnableInput component while editing.
 */

interface Args {
  model: Record<string, unknown>;
  renderGroup?: string;
  onChange?: (value: string) => void;
  onKeyUp?: (value: string) => void;
  modelValidations?: ValidationMap;
  groupName?: string;
}

export default class FormFieldGroupsComponent extends Component<Args> {
  @tracked showGroup: string | null = null;

  @action
  toggleGroup(group: string, isOpen: boolean) {
    this.showGroup = isOpen ? group : null;
  }

  get fieldGroups() {
    return this.args.groupName || 'fieldGroups';
  }
}
