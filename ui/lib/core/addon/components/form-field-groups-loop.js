import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module FormFieldGroupsLoop
 * FormFieldGroupsLoop components ARG TODO.
 *
 * @example
 * ```js
 <FormFieldGroupsLoop  ARG TODO/>
 * ```
 * @param {string} [mode=null] - ARG TODO
 */

// add group name to list here if you want to display within a flex box.
// check first no other group name in another model exists.
const MODEL_GROUPS_DISPLAY_FLEX = ['Key usage'];

export default class FormFieldGroupsLoop extends Component {
  @tracked flexGroups = [];
  constructor() {
    super(...arguments);
    let displayFlexGroups = this.args.model.fieldGroups.map((group) => {
      let key = Object.keys(group)[0]; // the key name e.g. default or Key usage
      return MODEL_GROUPS_DISPLAY_FLEX.includes(key) ? key : '';
    });
    console.log(displayFlexGroups, 'displayFlexGroups');
    this.flexGroups = displayFlexGroups;
  }
}
