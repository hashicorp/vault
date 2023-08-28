/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';
import IdentitySerializer from './_base';

export default IdentitySerializer.extend(EmbeddedRecordsMixin, {
  // we don't need to serialize relationships here
  serializeHasMany() {},
  attrs: {
    aliases: { embedded: 'always' },
  },
});
