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

  trueOrFalseString(model, field, returnTrue, returnFalse) {
    if (model.operationAll) {
      return returnTrue;
    }
    if (model.operationNone) {
      return returnFalse;
    }
    return model.get(field.name) ? returnTrue : returnFalse;
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
