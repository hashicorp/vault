import Component from '@ember/component';

export default Component.extend({
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
