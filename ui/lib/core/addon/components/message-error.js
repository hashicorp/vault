import { computed } from '@ember/object';
import Component from '@ember/component';
import layout from '../templates/components/message-error';

/**
 * @module MessageError
 * `MessageError` extracts an error from a model or a passed error and displays it using the `AlertBanner` component.
 *
 * @example
 * ```js
 * <MessageError @model={{model}} />
 * ```
 *
 * @param model=null{DS.Model} - An Ember data model that will be used to bind error statest to the internal
 * `errors` property.
 * @param errors=null{Array} - An array of error strings to show.
 * @param errorMessage=null{String} - An Error string to display.
 */
export default Component.extend({
  layout,
  model: null,
  errors: null,
  errorMessage: null,

  displayErrors: computed(
    'errorMessage',
    'model.{isError,adapterError.message,adapterError.errors.@each}',
    'errors',
    'errors.@each',
    function() {
      const errorMessage = this.get('errorMessage');
      const errors = this.get('errors');
      const modelIsError = this.get('model.isError');
      if (errorMessage) {
        return [errorMessage];
      }

      if (errors && errors.length > 0) {
        return errors;
      }

      if (modelIsError) {
        if (this.get('model.adapterError.errors.length') > 0) {
          return this.get('model.adapterError.errors');
        }
        return [this.get('model.adapterError.message')];
      }
    }
  ),
});
