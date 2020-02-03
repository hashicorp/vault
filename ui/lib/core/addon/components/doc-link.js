import Component from '@ember/component';
import { computed } from '@ember/object';
import hbs from 'htmlbars-inline-precompile';

/**
 * @module DocLink
 * `DocLink` components are used to render anchor links to relevant Vault documentation.
 *
 * @example
 * ```js
    <DocLink @path="/docs/secrets/kv/kv-v2.html">Learn about KV v2</DocLink>
 * ```
 *
 * @param path="/"{String} - The path to documentation on vaultproject.io that the component should link to.
 *
 */

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
