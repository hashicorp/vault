import TransformBase from './transform-edit-base';
import { computed } from '@ember/object';

export default TransformBase.extend({
  cliCommand: computed('model.{allowed_roles}', 'model.{type}', 'model.{tweak_source}', function() {
    if (!this.model) {
      return;
    }

    let { type, allowed_roles, tweak_source, name } = this.model;
    let wildCardRole = allowed_roles.find(role => role.includes('*'));

    // values to be returned
    let role = '<choose a role>';
    let value = 'value=<enter your value here>'; // change this when decode vs encode
    let tweak = '';

    // determine the role
    if (allowed_roles.length === 1 && !wildCardRole) {
      role = allowed_roles[0];
    }
    // determine the tweak_source
    if (type === 'fpe' && tweak_source === 'supplied') {
      tweak = 'tweak=<enter your tweak>';
    }

    return `${role} ${value} ${tweak} transformation=${name}`;
  }),
});
