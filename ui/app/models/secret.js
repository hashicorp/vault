import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import KeyMixin from 'vault/mixins/key-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend(KeyMixin, {
  failedServerRead: attr('boolean'),
  auth: attr('string'),
  lease_duration: attr('number'),
  lease_id: attr('string'),
  renewable: attr('boolean'),

  secretData: attr('object'),
  secretKeyAndValue: computed('secretData', function () {
    const data = this.secretData;
    return Object.keys(data).map((key) => {
      return { key, value: data[key] };
    });
  }),

  dataAsJSONString: computed('secretData', function () {
    return JSON.stringify(this.secretData, null, 2);
  }),

  isAdvancedFormat: computed('secretData', function () {
    const data = this.secretData;
    return data && Object.keys(data).some((key) => typeof data[key] !== 'string');
  }),

  helpText: attr('string'),
  // TODO this needs to be a relationship like `engine` on kv-v2
  backend: attr('string'),
  secretPath: lazyCapabilities(apiPath`${'backend'}/${'id'}`, 'backend', 'id'),
  canEdit: alias('secretPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  canRead: alias('secretPath.canRead'),
});
