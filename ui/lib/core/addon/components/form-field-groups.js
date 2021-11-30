import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/form-field-groups';

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
 *
 * @param [renderGroup=null] {String} - An allow list of groups to include in the render.
 * @param model=null {DS.Model} - Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered.
 * @param onChange=null {Func} - Handler that will get set on the `FormField` component.
 * @param onKeyUp=null {Func} - Handler that will set the value and trigger validation on input changes
 * @param validationMessages=null {Object} Object containing validation message for each property
 *
 */

export default Component.extend({
  layout,
  tagName: '',

  renderGroup: computed(function() {
    return null;
  }),

  model: null,

  onChange: () => {},
});
