import TransformBase from './transform-edit-base';

export default TransformBase.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'transform');
  },
});
