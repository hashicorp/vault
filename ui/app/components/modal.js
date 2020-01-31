/**
 * @module Modal
 * Modal components are used to...
 *
 * @example
 * ```js
 * <Modal @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';

export default Component.extend({
  title: null,
  onClose: () => {},
});
