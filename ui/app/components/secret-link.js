import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';
import { hrefTo } from 'vault/helpers/href-to';
const { computed } = Ember;

export function linkParams({ mode, secret, queryParams }) {
  let params;
  const route = `vault.cluster.secrets.backend.${mode}`;

  if (!secret || secret === ' ') {
    params = [route + '-root'];
  } else {
    params = [route, secret];
  }

  if (queryParams) {
    params.push(queryParams);
  }

  return params;
}

export default Ember.Component.extend({
  mode: 'list',

  secret: null,
  queryParams: null,
  ariaLabel: null,

  linkParams: computed('mode', 'secret', 'queryParams', function() {
    return linkParams(this.getProperties('mode', 'secret', 'queryParams'));
  }),

  attributeBindings: ['href', 'aria-label:ariaLabel'],

  href: computed('linkParams', function() {
    return hrefTo(this, ...this.get('linkParams'));
  }),

  layout: hbs`{{yield}}`,

  tagName: 'a',
});
