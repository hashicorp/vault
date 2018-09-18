'use strict';

module.exports = {
  extends: 'recommended',
  rules: {
    // these first 4 we should definitely make pass eventually
    // but they are likely their own PRs
    'no-partial': false,
    'simple-unless': false,
    'no-nested-interactive': false,
    'no-invalid-interactive': false,

    // if prettier ever gets this right I'd be in favor of enabling it
    'attribute-indentation': false,
    // I think the prettier glimmer parser is re-writing things incorrectly here
    // as I've fixed templates more than once - so leave this disabled
    'self-closing-void-elements': false,
  },
};
