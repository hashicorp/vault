/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module FormFieldGroups
 * `FormFieldGroups` components are field groups associated with a particular model. They render individual `FormField` components.
 *
 * @example
 * ```js
 * {{if model.fieldGroups}}
 *  <FormFieldGroups @model={{model}} />
 * {{/if}}
 *
 * ...
 *
 * <FormFieldGroups
 *  @model={{mountModel}}
 *  @onChange={{action "onTypeChange"}}
 *  @renderGroup="Method Options"
 *  @onKeyUp={{action "onKeyUp"}}
 *  @validationMessages={{validationMessages}}
 * />
 * ```
 * @callback onChangeCallback
 * @callback onKeyUpCallback
 * @param {Model} model- Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered.
 * @param {string} [renderGroup] - An allow list of groups to include in the render.
 * @param {onChangeCallback} [onChange] - Handler that will get set on the `FormField` component.
 * @param {onKeyUpCallback} [onKeyUp] - Handler that will set the value and trigger validation on input changes
 * @param {ModelValidations} [modelValidations] - Object containing validation message for each property
 * @param {string} [groupName='fieldGroups'] - attribute name where the field groups are
 */

export default class FormFieldGroupsComponent extends Component {
  @tracked showGroup = null;

  @action
  toggleGroup(group, isOpen) {
    this.showGroup = isOpen ? group : null;
  }

  get fieldGroups() {
    return this.args.groupName || 'fieldGroups';
  }
}
