/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { dateFromNow } from 'core/helpers/date-from-now';

export default class TransitFormShow extends Component {
  @service router;
  @service flashMessages;
  @service api;

  @action async rotateKey() {
    const { backend, id } = this.args.key;
    try {
      await this.api.secrets.transitRotateKey(id, backend, {});
      this.flashMessages.success('Key rotated.');
      // must refresh to see the updated versions, a model refresh does not trigger the change.
      await this.router.refresh();
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }

  // Symmetric key timestamps from the API are Unix seconds, not milliseconds.
  getTimestamp(time) {
    return dateFromNow([time * 1000], { addSuffix: true });
  }
}
