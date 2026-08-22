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
    const { backend, id } = this.args.form.data;
    try {
      await this.api.secrets.transitRotateKey(id, backend, {});
      this.flashMessages.success('Key rotated.');
      // must refresh to see the updated versions, a model refresh does not trigger the change.
      await this.router.refresh();
    } catch (e) {
      const { message } = this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }

  // Investigate - possibly a bug?
  // api returns the same value for creation time as ember data (eg. 1633024800) - but the date isn't rendering the same (ie. ED model pipe returns correct time but api value returns as 56 years ago)
  // not sure why the ED model data is taken as milliseconds and the api value is taken as seconds, but this is a workaround to get the correct time.
  getTimestamp(time) {
    return dateFromNow([time * 1000], { addSuffix: true });
  }
}
