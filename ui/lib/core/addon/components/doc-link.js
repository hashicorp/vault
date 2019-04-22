import Component from '@ember/component';
import { computed } from '@ember/object';
import hbs from 'htmlbars-inline-precompile';

export default Component.extend({
  tagName: 'a',
  classNames: ['doc-link'],
  attributeBindings: ['target', 'rel', 'href'],

  layout: hbs`{{yield}}`,

  target: '_blank',
  rel: 'noreferrer noopener',
  host: 'https://www.vaultproject.io',

  path: '/',
  href: computed('host', 'path', function() {
    return `${this.host}${this.path}`;
  }),
});
