import { computed } from '@ember/object';
import Mixin from '@ember/object/mixin';
import utils from '../lib/key-utils';

export default Mixin.create({
  flags: null,

  initialParentKey: null,

  isCreating: computed('initialParentKey', function() {
    return this.get('initialParentKey') != null;
  }),

  isFolder: computed('id', function() {
    return utils.keyIsFolder(this.get('id'));
  }),

  keyParts: computed('id', function() {
    return utils.keyPartsForKey(this.get('id'));
  }),

  parentKey: computed('id', 'isCreating', {
    get: function() {
      return this.get('isCreating') ? this.get('initialParentKey') : utils.parentKeyForKey(this.get('id'));
    },
    set: function(_, value) {
      return value;
    },
  }),

  keyWithoutParent: computed('id', 'parentKey', {
    get: function() {
      var key = this.get('id');
      return key ? key.replace(this.get('parentKey'), '') : null;
    },
    set: function(_, value) {
      if (value && value.trim()) {
        this.set('id', this.get('parentKey') + value);
      } else {
        this.set('id', null);
      }
      return value;
    },
  }),
});
