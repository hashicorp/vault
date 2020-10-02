/**
 * @module FormError
 * FormError components are used to...
 *
 * @example
 * ```js
 * <FormError @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
import layout from '../templates/components/form-error';

export default Component.extend({
  layout,
});
