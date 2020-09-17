import { belongsTo, attr } from '@ember-data/model';
import Secret from './secret';
import { computed } from '@ember/object';

export default Secret.extend({
  failedServerRead: attr('boolean'),
  pathAttr: 'path',
  version: attr('number'),
  secret: belongsTo('secret-v2'),
  path: attr('string'),
  deletionTime: attr('string'),
  createdTime: attr('string'),
  deleted: computed('deletionTime', function() {
    const deletionTime = new Date(this.deletionTime);
    const now = new Date();
    return deletionTime <= now;
  }),
  destroyed: attr('boolean'),
  currentVersion: attr('number'),
});
