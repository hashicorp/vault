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
import { inject as service } from '@ember/service';
import ShamirFlow from './shamir-flow';
import layout from '../templates/components/shamir-modal-flow';

export default ShamirFlow.extend({
  layout,
  store: service(),
  onClose: () => {},
  actions: {
    onCancelClose() {
      if (this.encoded_token) {
        this.send('reset');
      } else if (this.generateAction && !this.started) {
        if (this.generateStep !== 'chooseMethod') {
          this.send('reset');
        }
      } else {
        const adapter = this.get('store').adapterFor('cluster');
        adapter.generateDrOperationToken(this.model, { cancel: true });
        this.send('reset');
      }
      this.onClose();
    },
    onClose() {
      this.onClose();
    },
  },
});
