import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

const { Component, computed } = Ember;

export default Component.extend({
  tagName: 'a',
  attributeBindings: ['target', 'rel', 'href'],

  layout: hbs`{{yield}}`,

  target: '_blank',
  rel: 'noreferrer noopener',

  path: '/',
  href: computed('path', function() {
    return `https://www.vaultproject.io/docs${this.get('path')}`;
  }),
});
