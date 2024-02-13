/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { getOwner } from '@ember/application';
import config from '../config/environment';
import { service } from '@ember/service';

export default Component.extend({
  flashMessages: service(),
  auth: service(),

  useServiceWorker: null,

  async init() {
    this._super(...arguments);
    if (this.useServiceWorker === false) {
      return;
    }
    // check to see if we support ServiceWorker
    if ('serviceWorker' in navigator) {
      // this checks to see if there's an active service worker - if it failed to register
      // for any reason, then this would be null
      const worker = await navigator.serviceWorker.getRegistration(config.serviceWorkerScope);
      if (worker) {
        navigator.serviceWorker.addEventListener('message', this.serviceWorkerGetToken.bind(this));

        this.set('useServiceWorker', true);
      }
    }
  },
  willDestroy() {
    if (this.useServiceWorker) {
      navigator.serviceWorker.removeEventListener('message', this.serviceWorkerGetToken);
    }
    this._super(...arguments);
  },

  serviceWorkerGetToken(event) {
    const { action } = event.data;
    const [port] = event.ports;

    if (action === 'getToken') {
      port.postMessage({ token: this.auth.currentToken });
    } else {
      console.error('Unknown event', event); // eslint-disable-line
      port.postMessage({ error: 'Unknown request' });
    }
  },

  actions: {
    async removePeer(model) {
      const { nodeId } = model;
      try {
        await model.destroyRecord();
      } catch (e) {
        const errString = e.errors ? e.errors.join(' ') : e.message || e;
        this.flashMessages.danger(`There was an issue removing the peer ${nodeId}: ${errString}`);
        return;
      }
      this.flashMessages.success(`Successfully removed the peer: ${nodeId}.`);
    },

    downloadViaServiceWorker() {
      // the actual download happens when the user clicks the anchor link, and then the ServiceWorker
      // intercepts the request and adds auth headers.
      // Here we just want to notify users that something is happening before the browser starts the download
      this.flashMessages.success('The snapshot download will begin shortly.');
    },

    async downloadSnapshot() {
      // this entire method is the fallback behavior in case the browser either doesn't support ServiceWorker
      // or the UI is not being run on https.
      // here we're downloading the entire snapshot in memory, creating a dataurl with createObjectURL, and
      // then forcing a download by clicking a link that has a download attribute
      //
      // this is not the default because
      const adapter = getOwner(this).lookup('adapter:application');

      this.flashMessages.success('The snapshot download has begun.');
      let resp, blob;
      try {
        resp = await adapter.rawRequest('/v1/sys/storage/raft/snapshot', 'GET');
        blob = await resp.blob();
      } catch (e) {
        const errString = e.errors ? e.errors.join(' ') : e.message || e;
        this.flashMessages.danger(`There was an error trying to download the snapshot: ${errString}`);
      }
      const filename = 'snapshot.gz';
      const file = new Blob([blob], { type: 'application/x-gzip' });
      file.name = filename;
      if ('msSaveOrOpenBlob' in navigator) {
        navigator.msSaveOrOpenBlob(file, filename);
        return;
      }
      const a = document.createElement('a');
      const objectURL = window.URL.createObjectURL(file);
      a.href = objectURL;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(objectURL);
    },
  },
});
