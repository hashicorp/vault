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
 * @param model=null{DS.Model} - An Ember data model that will be used to bind error states to the internal
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
    'errors.[]',
    'model.adapterError.{errors.[],message}',
    'model.isError',
    'parentView.mountModel.{adapterError,isError}',
    function () {
      const errorMessage = this.errorMessage;
      const errors = this.errors;
      const modelIsError = this.model?.isError || this.parentView?.mountModel.isError;

      if (modelIsError) {
        let adapterError = this.model?.adapterError || this.parentView?.mountModel?.adapterError;
        if (!adapterError) {
          return;
        }
        if (adapterError.errors.length > 0) {
          return adapterError.errors.map((e) => {
            if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
            return e;
          });
        }
        return [adapterError.message];
      }
      if (errorMessage) {
        return [errorMessage];
      }

      if (errors && errors.length > 0) {
        return errors;
      }
      return 'no error';
    }
  ),
});
