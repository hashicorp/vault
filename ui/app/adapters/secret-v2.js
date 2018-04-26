import Ember from 'ember';
import SecretAdapter from './secret';

export default SecretAdapter.extend({
  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;

    return this.ajax(this.urlForSecret(snapshot.attr('backend'), id), 'POST', {
      data: { data },
    });
  },

  urlForSecret(backend, id, infix = 'data') {
    let url = `${this.buildURL()}/${backend}/${infix}/`;
    if (!Ember.isEmpty(id)) {
      url = url + id;
    }
    return url;
  },

  fetchByQuery(query, methodCall) {
    let { id, backend } = query;
    let args = [backend, id];
    if (methodCall === 'query') {
      args.push('metadata');
    }
    return this.ajax(this.urlForSecret(...args), 'GET', this.optionsForQuery(id, methodCall)).then(resp => {
      resp.id = id;
      return resp;
    });
  },
});
