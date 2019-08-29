import Component from '@ember/component';
import { getOwner } from '@ember/application';

export default Component.extend({
  actions: {
    async downloadSnapshot() {
      let adapter = getOwner(this).lookup('adapter:application');
      let resp = await adapter.rawRequest('/v1/sys/storage/raft/snapshot', 'GET');
      let blob = await resp.blob();

      let filename = 'raft.gzip';
      let file = new Blob([blob], { type: 'application/gzip' });
      file.name = filename;
      let a = document.createElement('a');
      a.href = window.URL.createObjectURL(file);
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
    },
  },
});
