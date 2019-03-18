import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module HomeLink
 * `HomeLink` is a span that contains either the text `home` or the `LogoEdition` component.
 *
 * @example
 * ```js
 * <HomeLink @class="navbar-item splash-page-logo">
 *  <LogoEdition />
 * </HomeLink>
 * ```
 *
 * @see {@link https://github.com/hashicorp/vault/search?l=Handlebars&q=HomeLink|Uses of HomeLink}
 * @see {@link https://github.com/hashicorp/vault/blob/master/ui/app/components/home-link.js|HomeLink Source Code}
 */

export default Component.extend({
  tagName: '',

  text: computed(function() {
    return 'home';
  }),

  computedClasses: computed('classNames', function() {
    return this.get('classNames').join(' ');
  }),
});
