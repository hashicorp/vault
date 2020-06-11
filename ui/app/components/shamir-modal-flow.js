/**
 * @module ShamirModalFlow
 * ShamirModalFlow components are used to...
 *
 * @example
 * ```js
 * <ShamirModalFlow @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import ShamirFlow from './shamir-flow';
import layout from '../templates/components/shamir-modal-flow';

export default ShamirFlow.extend({
  layout,
  onClose: () => {},
  actions: {
    onCancelClose() {
      console.log('on cancel close');
      this.reset();
      console.log('closing');
      this.onClose();
    },
  },
});
