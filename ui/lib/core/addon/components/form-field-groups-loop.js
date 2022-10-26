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
 */
/* eslint ember/no-empty-glimmer-component-classes: 'warn' */
export default class FormFieldGroupsLoop extends Component {}
