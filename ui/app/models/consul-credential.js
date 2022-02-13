import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const DISPLAY_FIELDS = ['accessor', 'token', 'local'];
export default Model.extend({
  helpText:
    'There are no inputs, just submit the form.',
  role: attr('object', {
    readOnly: true,
  }),

  local: attr('boolean'),
  accessor: attr('string'),
  token: attr('string'),

  attrs: computed('accessor', 'token', 'local', function() {
    if (this.accessKey || this.securityToken) {
      return expandAttributeMeta(this, DISPLAY_FIELDS.slice(0));
    }
    return expandAttributeMeta(this, []);
  }),

  toCreds: computed('accessor', 'token', 'local', function() {
    const props = {
      accessor: this.accessor,
      token: this.token,
      local: this.local
    };
    const propsWithVals = Object.keys(props).reduce((ret, prop) => {
      if (props[prop]) {
        ret[prop] = props[prop];
        return ret;
      }
      return ret;
    }, {});
    return JSON.stringify(propsWithVals, null, 2);
  }),
});
