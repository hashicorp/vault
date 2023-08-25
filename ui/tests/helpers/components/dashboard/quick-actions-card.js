const SELECTORS = {
  searchSelect: '.search-select',
  secretsEnginesSelect: '[data-test-secrets-engines-select]',
  actionSelect: '[data-test-select="action-select"]',
  emptyState: '[data-test-no-mount-selected-empty]',
  paramsTitle: '[data-test-search-select-params-title]',
  paramSelect: '[data-test-param-select]',
  getActionButton: (action) => `[data-test-button="${action}"]`,
};

export default SELECTORS;
