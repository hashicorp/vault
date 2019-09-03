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
    if ('serviceWorker' in navigator) {
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
      this.flashMessages.success('The snapshot download will begin shortly.');
    },

    async downloadSnapshot() {
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
      a.href = window.URL.createObjectURL(file);
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
    },
  },
});
