import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1/sys',
  pathForType(type) {
    let path = type.replace('policy', 'policies');
    return path;
  },

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor('policy');
    const data = serializer.serialize(snapshot);
    const name = snapshot.attr('name');

    return this.ajax(this.buildURL(type.modelName, name), 'PUT', { data }).then(() => {
      // doing this to make it like a Vault response - ember data doesn't like 204s if it's not a DELETE
      return {
        data: assign({}, snapshot.record.toJSON(), { id: name }),
      };
    });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments);
  },

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> add policy search select and group search select to group creation
  query(store, type) {
    return this.ajax(this.buildURL(type.modelName), 'GET', {
      data: { list: true },
    });
<<<<<<< HEAD
=======
  query(store, type, options) {
    return this.ajax(
      options && options.findAll ? this.buildURL() + `/${type.modelName}` : this.buildURL(type.modelName),
      'GET',
      {
        data: { list: true },
      }
    );
>>>>>>> add and use ember-power-select for policies in entity create/edit
=======
>>>>>>> add policy search select and group search select to group creation
  },
});
