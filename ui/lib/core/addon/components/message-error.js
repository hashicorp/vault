import Component from '@glimmer/component';
import layout from '../templates/components/message-error';
import { setComponentTemplate } from '@ember/component';

/**
 * @module MessageError
 * `MessageError` extracts an error from a model or a passed error and displays it using the `AlertBanner` component.
 *
 * @example
 * ```js
 * <MessageError @model={{model}} />
 * ```
 *
 * @param {object} [model=null] - An Ember data model that will be used to bind error states to the internal
 * `errors` property.
 * @param {array} [errors=null] - An array of error strings to show.
 * @param {string} [errorMessage=null] - An Error string to display.
 */

class MessageError extends Component {
  get errorMessage() {
    return this.args.errorMessage;
  }

  get errors() {
    return this.args.errors;
  }

  get model() {
    return this.args.model;
  }

  get displayErrors() {
    if (this.errorMessage) {
      return [this.errorMessage];
    }

    if (this.errors && this.errors.length > 0) {
      return this.errors;
    }

    if (this.model?.isError) {
      let adapterError = this.model?.adapterError;
      if (!adapterError) {
        return null;
      }
      if (adapterError.errors.length > 0) {
        return adapterError.errors.map((e) => {
          if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
          return e;
        });
      }
      return [adapterError.message];
    }
    return null;
  }
}
export default setComponentTemplate(layout, MessageError);

// import { computed } from '@ember/object';
// import Component from '@ember/component';
// import layout from '../templates/components/message-error';

// /**
//  * @module MessageError
// 	@@ -11,48 +11,39 @@ import layout from '../templates/components/message-error';
//  * <MessageError @model={{model}} />
//  * ```
//  *
//  * @param model=null{DS.Model} - An Ember data model that will be used to bind error statest to the internal
//  * `errors` property.
//  * @param errors=null{Array} - An array of error strings to show.
//  * @param errorMessage=null{String} - An Error string to display.
//  */
// export default Component.extend({
//   layout,
//   model: null,
//   errors: null,
//   errorMessage: null,

//   displayErrors: computed(
//     'errorMessage',
//     'model.{isError,adapterError.message,adapterError.errors.@each}',
//     'errors',
//     'errors.[]',
//     function () {
//       const errorMessage = this.errorMessage;
//       const errors = this.errors;
//       const modelIsError = this.model?.isError;
//       if (errorMessage) {
//         return [errorMessage];
//       }

//       if (errors && errors.length > 0) {
//         return errors;
//       }

//       if (modelIsError) {
//         if (!this.model.adapterError) {
//           return;
//         }
//         if (this.model.adapterError.errors.length > 0) {
//           return this.model.adapterError.errors.map((e) => {
//             if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
//             return e;
//           });
//         }
//         return [this.model.adapterError.message];
//       }

//       return 'no error';
//     }
//   ),
// });
