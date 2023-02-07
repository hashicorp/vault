import IdentityAdapter from './base';

export default IdentityAdapter.extend({
  lookup(store, data) {
    const url = `/${this.urlPrefix()}/identity/lookup/group`;
    return this.ajax(url, 'POST', { data }).then((response) => {
      // unsuccessful lookup is a 204
      if (!response) return;
      const modelName = 'identity/group';
      store.push(
        store
          .serializerFor(modelName)
          .normalizeResponse(store, store.modelFor(modelName), response, response.data.id, 'findRecord')
      );
      return response;
    });
  },
});
