/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { getOwner } from '@ember/application';
import { service } from '@ember/service';
import { alias } from '@ember/object/computed';
import { AbortController } from 'fetch';

export default Component.extend({
  file: null,
  errors: null,
  forceRestore: false,
  flashMessages: service(),
  isUploading: alias('restore.isRunning'),
  abortController: null,
  restore: task(function* () {
    this.set('errors', null);
    const adapter = getOwner(this).lookup('adapter:application');
    try {
      let url = '/v1/sys/storage/raft/snapshot';
      if (this.forceRestore) {
        url = `${url}-force`;
      }
      const file = new Blob([this.file], { type: 'application/gzip' });
      const controller = new AbortController();
      this.set('abortController', controller);
      yield adapter.rawRequest(url, 'POST', { body: file, signal: controller.signal });
      this.flashMessages.success('The snapshot was successfully uploaded!');
    } catch (e) {
      if (e.name === 'AbortError') {
        return;
      }
      let resp;
      if (e.json) {
        resp = yield e.json();
      }
      const err = resp ? resp.errors : [e];
      this.set('errors', err);
    }
  }),
  actions: {
    cancelUpload() {
      this.abortController.abort();
    },
  },
});
