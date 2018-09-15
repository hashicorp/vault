import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',

  /*
   * @public String
   * A whitelist of groups to include in the render
   */
  renderGroup: computed(function() {
    return null;
  }),

  /*
   * @public DS.Model
   * model to be passed down to form-field component
   * if `fieldGroups` is present on the model then it will be iterated over and
   * groups of `form-field` components will be rendered
   *
   */
  model: null,

  /*
   * @public Function
   * onChange handler that will get set on the form-field component
   *
   */
  onChange: () => {},
});
