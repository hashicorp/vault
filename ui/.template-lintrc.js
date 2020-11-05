'use strict';

module.exports = {
  extends: 'recommended',
  rules: {
    // should definitely move to template only
    // glimmer components for this one
    'no-partial': false,

    // these need to be looked into, but
    // may be a bigger change
    'no-invalid-interactive': false,
    'simple-unless': false,

    'self-closing-void-elements': false,
    'no-unnecessary-concat': false,
    'no-quoteless-attributes': false,
    'no-nested-interactive': false,

    // not sure we'll ever want these on,
    // would be nice but if prettier isn't doing
    // it for us, then not sure it's worth it
    'attribute-indentation': false,
    'block-indentation': false,
    quotes: false,
  },
};
