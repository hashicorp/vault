import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { getOwner } from '@ember/application';

export default Component.extend({
  file: null,
  errors: null,
  forceRestore: false,
  restore: task(function*() {
    this.set('errors', null);
    let adapter = getOwner(this).lookup('adapter:application');
    try {
      let url = '/v1/sys/storage/raft/snapshot';
      if (this.forceRestore) {
        url = `${url}-force`;
      }
      let file = new Blob([this.file], { type: 'application/gzip' });
      yield adapter.rawRequest(url, 'POST', { body: file });
    } catch (e) {
      let resp;
      if (e.json) {
        resp = yield e.json();
      }
      let err = resp ? resp.errors : [e];
      this.set('errors', err);
    }
  }),
});
