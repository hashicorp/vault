import Ember from 'ember';

export default Ember.Component.extend({
  /*
   * @public
   * @param DS.Model
   *
   * the pki-certificate model
   */
  item: null,

  actions: {
    delete(item) {
      item.save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
