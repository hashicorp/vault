import { clickable, fillable } from 'ember-cli-page-object';

import fields from '../form-field';
export default {
  ...fields,
  ttlValue: fillable('[data-test-ttl-value]'),
  ttlUnit: fillable('[data-test-ttl-value]'),
  save: clickable('[data-test-save-config]'),
};
