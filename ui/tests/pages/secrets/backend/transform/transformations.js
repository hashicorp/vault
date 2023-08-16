import { create, clickable, fillable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list'),
  visitShow: visitable('/vault/secrets/:backend/show/:id'),
  visitCreate: visitable('/vault/secrets/:backend/create'),
  createLink: clickable('[data-test-secret-create]'),
  name: fillable('[data-test-input="name"]'),
  submit: clickable('[data-test-transform-create]'),
  type: fillable('[data-test-input="type"'),
  tweakSource: fillable('[data-test-input="tweak_source"'),
  maskingChar: fillable('[data-test-input="masking_character"'),
  save: clickable('[data-test-transformation-save-button]'),
});
