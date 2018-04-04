import Ember from 'ember';

const { get, set } = Ember;

export default Ember.Component.extend({
  config: null,
  mounts: null,

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: ['cubbyhole', 'system', 'token', 'identity'],

  actions: {
    addOrRemovePath(path, e) {
      let config = get(this, 'config') || [];
      let paths = get(config, 'paths').slice();

      if (e.target.checked) {
        paths.addObject(path);
      } else {
        paths.removeObject(path);
      }

      set(config, 'paths', paths);
    },
  },
});
