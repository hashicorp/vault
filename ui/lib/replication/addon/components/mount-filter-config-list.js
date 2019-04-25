import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import layout from '../templates/components/mount-filter-config-list';

export default Component.extend({
  layout,
  config: null,
  mounts: null,

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: computed(function() {
    return ['cubbyhole', 'system', 'token', 'identity', 'ns_system', 'ns_identity'];
  }),

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
