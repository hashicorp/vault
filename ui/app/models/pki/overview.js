import Model, { attr } from '@ember-data/model';
// import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
// import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class PkiOverview extends Model {
  @attr('string', {
    label: 'Label of a form field',
    subText: 'Subtext of a form field.',
  })
  name;

  @attr('string')
  backend;

  icon = 'key';
  // follow the key model from keymgt
}
