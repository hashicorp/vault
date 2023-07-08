/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module FormFieldGroupsLoop
 * FormFieldGroupsLoop components loop through the "groups" set on a model and display them either as default or behind toggle components.
 *
 * @example
 * ```js
 <FormFieldGroupsLoop @model={{this.model}} @mode={{if @model.isNew "create" "update"}}/>
 * ```
 * @param {class} model - The routes model class.
 * @param {string} mode - "create" or "update" used to hide the name form field. TODO: not ideal, would prefer to disable it to follow new design patterns.
 * @param {function} [modelValidations] - Passed through to formField.
 * @param {boolean} [showHelpText] - Passed through to formField.
 * @param {string} [groupName="fieldGroups"] - option to override key on the model where groups are located
 */
export default class FormFieldGroupsLoop extends Component {
  get fieldGroups() {
    return this.args.groupName || 'fieldGroups';
  }
}
