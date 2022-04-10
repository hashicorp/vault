import Component from '@ember/component';

export default Component.extend({
  init() {
    this._super(...arguments);
    this.set('isOpen', this.customPath !== '');
  },
});
