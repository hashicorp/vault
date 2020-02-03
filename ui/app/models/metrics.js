import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

let { attr } = DS;

export default DS.Model.extend({
  // console.log("getting hit")
});
