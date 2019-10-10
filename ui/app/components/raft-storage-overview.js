import Component from '@ember/component';
import { getOwner } from '@ember/application';
import config from '../config/environment';
import { inject as service } from '@ember/service';

export default Component.extend({
  flashMessages: service(),
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
      let worker = await navigator.serviceWorker.getRegistration(config.serviceWorkerScope);
      if (worker) {
        this.set('useServiceWorker', true);
      }
    }
  },

  actions: {
    async removePeer(model) {
      let { nodeId } = model;
      try {
        await model.destroyRecord();
      } catch (e) {
        let errString = e.errors ? e.errors.join(' ') : e.message || e;
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
      let adapter = getOwner(this).lookup('adapter:application');

      this.flashMessages.success('The snapshot download has begun.');
      let resp, blob;
      try {
        resp = await adapter.rawRequest('/v1/sys/storage/raft/snapshot', 'GET');
        blob = await resp.blob();
      } catch (e) {
        let errString = e.errors ? e.errors.join(' ') : e.message || e;
        this.flashMessages.danger(`There was an error trying to download the snapshot: ${errString}`);
      }
      let filename = 'snapshot.gz';
      let file = new Blob([blob], { type: 'application/x-gzip' });
      file.name = filename;
      if ('msSaveOrOpenBlob' in navigator) {
        navigator.msSaveOrOpenBlob(file, filename);
        return;
      }
      let a = document.createElement('a');
      let objectURL = window.URL.createObjectURL(file);
      a.href = objectURL;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(objectURL);
    },
  },
});
