/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { setIdentifierGenerationMethod } from '@ember-data/store';
import { dasherize } from '@ember/string';
import { v4 as uuidv4 } from 'uuid';

export function initialize() {
  // see this GH issue for more information https://github.com/emberjs/data/issues/8106
  // Ember Data uses uuidv4 library to generate ids which relies on the crypto API which is no available in unsecure contexts
  // the suggested polyfill was added in 4.6.2 so until we upgrade we need to define our own id generation method
  // https://api.emberjs.com/ember-data/4.5/classes/IdentifierCache/methods/getOrCreateRecordIdentifier?anchor=getOrCreateRecordIdentifier
  // the uuid library was brought in to replace other usages of crypto in the app so it is safe to use in unsecure contexts
  // adapted from defaultGenerationMethod -- https://github.com/emberjs/data/blob/v4.5.0/packages/store/addon/-private/identifier-cache.ts#LL82-L94C2
  setIdentifierGenerationMethod((data) => {
    if (data.lid) {
      return data.lid;
    }
    if (data.id) {
      return `@lid:${dasherize(data.type)}-${data.id}`;
    }
    return uuidv4();
  });
}

export default {
  name: 'ember-data-identifiers',
  initialize,
};
