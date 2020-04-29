import Component from '@ember/component';
import layout from '../templates/components/replication-header';

export default Component.extend({
  layout,
  isSecondary: null,
  'data-test-replication-header': true,
});
