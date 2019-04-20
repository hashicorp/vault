import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import layout from '../templates/components/replication-summary';

export default Component.extend({
  layout,
  showModeSummary: false,
  cluster: null,
  replicationAttrs: alias('cluster.replicationAttrs'),
});
