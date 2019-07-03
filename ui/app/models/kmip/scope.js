import { computed } from '@ember/object';
import DS from 'ember-data';
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr } = DS;
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

let Model = DS.Model.extend({
  name: attr('string'),
  backend: attr({ readOnly: true }),
  attrs: computed(function() {
    return expandAttributeMeta(this, ['name']);
  }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/scope/${'id'}`,
});
