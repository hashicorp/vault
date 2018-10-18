import Secret from './secret';
import DS from 'ember-data';
import { bool } from '@ember/object/computed';

const { attr, belongsTo } = DS;

export default Secret.extend({
  pathAttr: 'path',
  version: attr('number'),
  secret: belongsTo('secret-v2'),
  path: attr('string'),
  deletionTime: attr('string'),
  createdTime: attr('string'),
  deleted: bool('deletionTime'),
  destroyed: attr('boolean'),
  currentVersion: attr('number'),
});
