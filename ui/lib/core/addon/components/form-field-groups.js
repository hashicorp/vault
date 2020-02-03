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
 * />
 * ```
 *
 * @param [renderGroup=null] {String} - A whitelist of groups to include in the render.
 * @param model=null {DS.Model} - Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered.
 * @param onChange=null {Func} - Handler that will get set on the `FormField` component.
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
