import Ember from 'ember';
import utils from '../lib/key-utils';

export default Ember.Mixin.create({
  flags: null,

  initialParentKey: null,

  isCreating: Ember.computed('initialParentKey', function() {
    return this.get('initialParentKey') != null;
  }),

  isFolder: Ember.computed('id', function() {
    return utils.keyIsFolder(this.get('id'));
  }),

  keyParts: Ember.computed('id', function() {
    return utils.keyPartsForKey(this.get('id'));
  }),

  parentKey: Ember.computed('id', 'isCreating', {
    get: function() {
      return this.get('isCreating') ? this.get('initialParentKey') : utils.parentKeyForKey(this.get('id'));
    },
    set: function(_, value) {
      return value;
    },
  }),

  keyWithoutParent: Ember.computed('id', 'parentKey', {
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
