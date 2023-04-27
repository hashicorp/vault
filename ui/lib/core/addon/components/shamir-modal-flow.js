/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module ShamirModalFlow
 * ShamirModalFlow is an extension of the ShamirFlow component that does the Generate Action Token workflow inside of a Modal.
 * Please note, this is not an extensive list of the required parameters -- please see ShamirFlow for others
 *
 * @example
 * ```js
 * <ShamirModalFlow @onClose={action 'onClose'}>This copy is the main paragraph when the token flow has not started</ShamirModalFlow>
 * ```
 * @param {function} onClose - This function will be triggered when the modal intends to be closed
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
        const adapter = this.store.adapterFor('cluster');
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
