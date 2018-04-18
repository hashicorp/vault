import Ember from 'ember';

// removes Ember Data records from the cache when the model
// changes or you move away from the current route
export default Ember.Mixin.create({
  modelPath: 'model',
  unloadModel() {
    const model = this.controller.get(this.get('modelPath'));
    if (!model || !model.unloadRecord) {
      return;
    }
    this.store.unloadRecord(model);
    model.destroy();
  },

  actions: {
    willTransition() {
      this.unloadModel();
      return true;
    },
  },
});
