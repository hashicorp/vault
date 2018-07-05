import Ember from 'ember';
import DS from 'ember-data';
const { attr, belongsTo, hasMany } = DS;
const { computed } = Ember;

import { memberAction } from 'ember-api-actions';

export default DS.Model.extend({
  approved: attr('boolean'),
  requestPath: attr('string'),

  requestEntity: belongsTo('identity/entity', { async: false }),
  authorizations: hasMany('identity/entity', { async: false }),

  request: memberAction({
    path: 'request',
    type: 'post',
    urlType: 'queryRecord',
  }),

  authorize: memberAction({
    path: 'authorize',
    type: 'post',
    urlType: 'queryRecord',
  }),
});
