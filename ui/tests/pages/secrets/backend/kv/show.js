import { Base } from '../show';
import { create, clickable, collection, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  rows: collection({
    scope: 'data-test-row-label',
  }),
  edit: clickable('[data-test-secret-json-toggle]'),
  editIsPresent: isPresent('[data-test-secret-json-toggle]'),
});
