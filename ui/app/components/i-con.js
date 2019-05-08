import { camelize } from '@ember/string';
import Component from '@ember/component';
import { computed } from '@ember/object';
import hbs from 'htmlbars-inline-precompile';

/**
 * @module ICon
 * `ICon` components are glyphs used to indicate important information.
 *
 * @example
 * ```js
 * <ICon @glyph="cancel-square-outline" />
 * ```
 * @param glyph=null {String} - The glyph type.
 *
 */

export const GLYPHS_WITH_SVG_TAG = [
  'cancel-square-outline',
  'cancel-square-fill',
  'check-circle-fill',
  'check-plain',
  'checkmark-circled-outline',
  'close-circled-outline',
  'console',
  'control-lock',
  'docs',
  'download',
  'edition-enterprise',
  'edition-oss',
  'false',
  'file',
  'folder',
  'hidden',
  'information-reversed',
  'learn',
  'neutral-circled-outline',
  'perf-replication',
  'person',
  'role',
  'status-indicator',
  'stopwatch',
  'tour',
  'true',
  'upload',
  'video',
  'visible',
];

export default Component.extend({
  layout: hbs`
    {{#if excludeSVG}}
      {{partial partialName}}
    {{else}}
      <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="{{size}}" height="{{size}}" viewBox="0 0 512 512">
        {{partial partialName}}
      </svg>
    {{/if}}
  `,

  tagName: 'span',
  excludeIconClass: false,
  classNameBindings: ['excludeIconClass::icon'],
  classNames: ['has-current-color-fill'],

  attributeBindings: ['aria-label', 'aria-hidden'],

  glyph: null,

  excludeSVG: computed('glyph', function() {
    let glyph = this.get('glyph');
    return glyph.startsWith('enable/') || GLYPHS_WITH_SVG_TAG.includes(glyph);
  }),

  size: computed('glyph', function() {
    return this.get('glyph').startsWith('enable/') ? 48 : 12;
  }),

  partialName: computed('glyph', function() {
    const glyph = this.get('glyph');
    return `svg/icons/${camelize(glyph)}`;
  }),
});
