import { create, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  metrics: visitable('/vault/metrics'),
  config: visitable('/vault/metrics/config'),
  configEdit: visitable('/vault/metrics/edit'),
  startInput: fillable('[data-test-start-input]'),
  endInput: fillable('[data-test-end-input]'),
  queryButton: clickable('[data-test-metrics-query]'),
  configTab: clickable('[data-test-configuration-tab]'),
  metricsTab: clickable('[data-test-usage-tab]'),
});
