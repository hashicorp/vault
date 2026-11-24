/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { alias } from '@ember/object/computed';

export default Component.extend({
  file: null,
  errors: null,
  forceRestore: false,
  api: service(),
  flashMessages: service(),
  isUploading: alias('restore.isRunning'),
  abortController: null,
  restore: task(function* () {
    this.set('errors', null);
    try {
      let url = '/sys/storage/raft/snapshot';
      if (this.forceRestore) {
        url = `${url}-force`;
      }
      const file = new Blob([this.file], { type: 'application/gzip' });
      const controller = new AbortController();
      this.set('abortController', controller);
      // the snapshot is a binary upload with no generated API method
      yield this.api.request.post(url, file, undefined, controller.signal);
      this.flashMessages.success('The snapshot was successfully uploaded!');
    } catch (e) {
      if (e.name === 'AbortError') {
        return;
      }
      const { message, response } = yield this.api.parseError(e);
      this.set('errors', response?.errors || [message]);
    }
  }),
  actions: {
    cancelUpload() {
      this.abortController.abort();
    },
  },
});
