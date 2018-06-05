import ApplicationAdapter from './application';
import DS from 'ember-data';
import Ember from 'ember';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'capabilities-self';
  },

  findRecord(store, type, id) {
    return this.ajax(this.buildURL(type), 'POST', { data: { paths: [id] } }).catch(e => {
      if (e instanceof DS.AdapterError) {
        Ember.set(e, 'policyPath', 'sys/capabilities-self');
      }
      throw e;
    });
  },

  queryRecord(store, type, query) {
    const { id } = query;
    return this.findRecord(store, type, id).then(resp => {
      resp.path = id;
      return resp;
    });
  },
});
