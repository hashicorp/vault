import { computed } from '@ember/object';
import utils from 'vault/lib/key-utils';

export function withKeyMixin() {
  return function decorator(SuperClass) {
    return class KeyMixin extends SuperClass {
      // what attribute has the path for the key
      // will.be 'path' for v2 or 'id' v1
      pathAttr = 'path';
      flags = null;
      initialParentKey = null;

      @computed('initialParentKey')
      get isCreating() {
        return this.initialParentKey != null;
      }

      pathVal() {
        return this[this.pathAttr] || this.id;
      }

      // rather than using defineProperty for all of these,
      // we're just going to hardcode the known keys for the path ('id' and 'path')
      @computed('id', 'path')
      get isFolder() {
        return utils.keyIsFolder(this.pathVal());
      }

      @computed('id', 'path')
      get keyParts() {
        return utils.keyPartsForKey(this.pathVal());
      }

      @computed('id', 'path', 'isCreating', 'initialParentKey')
      get parentKey() {
        return this.isCreating ? this.initialParentKey : utils.parentKeyForKey(this.pathVal());
      }
      set parentKey(value) {
        this.parentKey = value;
      }

      @computed('id', 'path', 'parentKey')
      get keyWithoutParent() {
        var key = this.pathVal();
        return key ? key.replace(this.parentKey, '') : null;
      }
      set keyWithoutParent(value) {
        if (value && value.trim()) {
          this[this.pathAttr] = this.parentKey + value;
        } else {
          this[this.pathAttr] = null;
        }
        this.keyWithoutParent = value;
      }
    };
  };
}
