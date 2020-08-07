import Secret from './secret';
import DS from 'ember-data';
import { computed } from '@ember/object';

const { attr, belongsTo } = DS;

export default Secret.extend({
  failedServerRead: attr('boolean'),
  pathAttr: 'path',
  version: attr('number'),
  secret: belongsTo('secret-v2'),
  path: attr('string'),
  deletionTime: attr('string'),
  createdTime: attr('string'),
  deleted: computed('deletionTime', function() {
    const deletionTime = new Date(this.get('deletionTime'));
    const now = new Date();
    return deletionTime <= now;
  }),
  destroyed: attr('boolean'),
  currentVersion: attr('number'),
});
