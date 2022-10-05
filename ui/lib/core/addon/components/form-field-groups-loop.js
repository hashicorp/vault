import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module FormFieldGroupsLoop
 * FormFieldGroupsLoop components loop through the "groups" set on a model and display them either as default or behind toggle components. 
 * Option to have toggle components display the fields within a css grid. 
 * See the "Key usage" section on the PKI engine for an example.
 * To use css grid, you must add the group name to the array MODEL_GROUPS_DISPLAY_GRID.
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

// Add group name to list here if you want to display within a css grid.
// Check first that no other group name exists in another model.
const MODEL_GROUPS_DISPLAY_GRID = ['Key usage'];

export default class FormFieldGroupsLoop extends Component {
  @tracked gridGroups = [];
  constructor() {
    super(...arguments);
    let displayGridGroups = this.args.model.fieldGroups?.filter((group) => {
      let key = Object.keys(group)[0]; // ex: 'Key usage' or 'default'
      return MODEL_GROUPS_DISPLAY_GRID.includes(key);
    });
    if (displayGridGroups.length === 0) {
      return;
    }
    this.gridGroups = Object.keys(displayGridGroups[0]);
  }
}
