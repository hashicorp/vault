import Secret from './secret';
import DS from 'ember-data';

const { attr, belongsTo } = DS;

export default Secret.extend({
  pathAttr: 'path',
  version: attr('number'),
  secret: belongsTo('secret-v2'),
  path: attr('string'),
  currentVersion: attr('number'),
});
