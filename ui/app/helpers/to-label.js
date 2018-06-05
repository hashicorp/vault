import Ember from 'ember';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';

export function toLabel(val) {
  return capitalize([humanize([dasherize(val)])]);
}

export default Ember.Helper.helper(toLabel);
