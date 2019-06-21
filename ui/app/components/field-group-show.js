/**
 * @module FieldGroupShow
 * FieldGroupShow components loop through a Model's fieldGroups
 * to display their attributes
 *
 * @example
 * ```js
 * <FieldGroupShow @model={{model}} @showAllFields=true />
 * ```
 *
 * @param model {Object} - the model
 * @param [showAllFields=false] {boolean} - whether to show fields with empty values
 */
import Component from '@ember/component';
import layout from '../templates/components/field-group-show';

export default Component.extend({
  layout,
  model: null,
  showAllFields: false,
});
