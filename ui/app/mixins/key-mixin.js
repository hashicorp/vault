/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { computed } from '@ember/object';
import Mixin from '@ember/object/mixin';
import { keyIsFolder, keyPartsForKey, parentKeyForKey } from 'core/utils/key-utils';

export default Mixin.create({
  // what attribute has the path for the key
  // will.be 'path' for v2 or 'id' v1
  pathAttr: 'path',
  flags: null,

  initialParentKey: null,

  isCreating: computed('initialParentKey', function () {
    return this.initialParentKey != null;
  }),

  pathVal() {
    return this[this.pathAttr] || this.id;
  },

  // rather than using defineProperty for all of these,
  // we're just going to hardcode the known keys for the path ('id' and 'path')
  isFolder: computed('id', 'path', function () {
    return keyIsFolder(this.pathVal());
  }),

  keyParts: computed('id', 'path', function () {
    return keyPartsForKey(this.pathVal());
  }),

  parentKey: computed('id', 'path', 'isCreating', {
    get: function () {
      return this.isCreating ? this.initialParentKey : parentKeyForKey(this.pathVal());
    },
    set: function (_, value) {
      return value;
    },
  }),

  keyWithoutParent: computed('id', 'path', 'parentKey', {
    get: function () {
      var key = this.pathVal();
      return key ? key.replace(this.parentKey, '') : null;
    },
    set: function (_, value) {
      if (value && value.trim()) {
        this.set(this.pathAttr, this.parentKey + value);
      } else {
        this.set(this.pathAttr, null);
      }
      return value;
    },
  }),
});
