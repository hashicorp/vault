/**
 * @module OperationFieldDisplay
 * OperationFieldDisplay components are used on KMIP role show pages to display the allowed operations on that model
 *
 * @example
 * ```js
 * <OperationFieldDisplay @model={{model}} />
 * ```
 *
 * @param model {DS.Model} - model is the KMIP role model that needs to display its allowed operations
 *
 */
import Component from '@ember/component';
import layout from '../templates/components/operation-field-display';

export default Component.extend({
  layout,
  tagName: '',
  model: null,

  trueOrFalseString(model, field, trueString, falseString) {
    if (model.operationAll) {
      return trueString;
    }
    if (model.operationNone) {
      return falseString;
    }
    return model.get(field.name) ? trueString : falseString;
  },

  actions: {
    iconClass(model, field) {
      return this.trueOrFalseString(model, field, 'icon-true', 'icon-false');
    },
    iconGlyph(model, field) {
      return this.trueOrFalseString(model, field, 'check-circle-outline', 'cancel-square-outline');
    },
  },
});
