import ApplicationAdapter from './application';
import Ember from 'ember';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  defaultSerializer: 'ssh',

  url(/*role*/) {
    Ember.assert('Override the `url` method to extend the SSH adapter', false);
  },

  createRecord(store, type, snapshot, requestType) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const role = snapshot.attr('role');

    return this.ajax(this.url(role, snapshot), 'POST', { data }).then(response => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
