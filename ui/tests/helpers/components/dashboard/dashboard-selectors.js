export const SELECTORS = {
  cardName: (name) => `[data-test-card="${name}"]`,
  emptyState: (name) => `[data-test-empty-state="${name}"]`,
  cardHeader: (name) => `[data-test-dashboard-card-header="${name}"]`,
  tableRow: (name) => `[data-test-dashboard-table="${name}"] tr`,
};
