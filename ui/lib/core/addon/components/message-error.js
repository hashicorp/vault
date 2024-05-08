/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import layout from '../templates/components/message-error';
import { setComponentTemplate } from '@ember/component';

/**
 * @module MessageError
 * Renders form errors using the <Hds::Alert> component and extracts errors from
 * a model, passed errorMessage or array of errors and displays each in a separate banner.
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
  get displayErrors() {
    const { errorMessage, errors, model } = this.args;
    if (errorMessage) {
      return [errorMessage];
    }

    if (errors && errors.length > 0) {
      return errors;
    }

    if (model?.isError) {
      const adapterError = model?.adapterError;
      if (!adapterError) {
        return null;
      }
      if (adapterError.errors?.length > 0) {
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
