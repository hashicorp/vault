import { set } from '@ember/object';
import ApplicationAdapter from './application';
import DS from 'ember-data';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'capabilities-self';
  },

  findRecord(store, type, id) {
    return this.ajax(this.buildURL(type), 'POST', { data: { paths: [id] } }).catch(e => {
      if (e instanceof DS.AdapterError) {
        set(e, 'policyPath', 'sys/capabilities-self');
      }
      throw e;
    });
  },

  queryRecord(store, type, query) {
    const { id } = query;
    if (!id) {
      return;
    }
    return this.findRecord(store, type, id).then(resp => {
      resp.path = id;
      return resp;
    });
  },
});
