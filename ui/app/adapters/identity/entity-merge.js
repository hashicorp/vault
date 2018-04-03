import IdentityAdapter from './base';

export default IdentityAdapter.extend({
  buildURL() {
    // first arg is modelName which we're hardcoding in the call to _super.
    let [, ...args] = arguments;
    return this._super('identity/entity/merge', ...args);
  },

  createRecord(store, type, snapshot) {
    return this._super(...arguments).then(() => {
      // return the `to` id here so we can redirect to it on success
      // (and because ember _loves_ 204s for createRecord)
      return { id: snapshot.attr('toEntityId') };
    });
  },
});
