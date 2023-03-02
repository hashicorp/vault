import { setIdentifierGenerationMethod } from '@ember-data/store';
import { v4 as uuidv4 } from 'uuid';

export function initialize() {
  // see this GH issue for more information https://github.com/emberjs/data/issues/8106
  // Ember Data uses uuidv4 library to generate ids which relies on the crypto API which is no available in unsecure contexts
  // the suggested polyfill was added in 4.6.2 so until we upgrade we need to define our own id generation method
  // the uuid library was brought in to replace other usages of crypto in the app so it is safe to use in unsecure contexts
  setIdentifierGenerationMethod((resource) => {
    return resource.lid || uuidv4();
  });
}

export default {
  name: 'ember-data-identifiers',
  initialize,
};
