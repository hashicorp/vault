import Mixin from '@ember/object/mixin';

// removes Ember Data records from the cache when the model
// changes or you move away from the current route
export default Mixin.create({
  modelPath: 'model',
  unloadModel() {
    let { modelPath } = this;
    let model = this.controller.get(modelPath);
    if (!model || !model.unloadRecord) {
      return;
    }
    this.store.unloadRecord(model);
    model.destroy();
    // it's important to unset the model on the controller since controllers are singletons
    this.controller.set(modelPath, null);
  },

  actions: {
    willTransition() {
      this.unloadModel();
      return true;
    },
  },
});
