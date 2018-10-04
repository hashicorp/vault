import Secret from './secret';
import DS from 'ember-data';

const { attr } = DS;

export default Secret.extend({
  version: attr('number'),
});
