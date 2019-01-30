import { camelize } from '@ember/string';
import Component from '@ember/component';
import { computed } from '@ember/object';
import hbs from 'htmlbars-inline-precompile';

const GLYPHS_WITH_SVG_TAG = [
  'learn',
  'video',
  'tour',
  'stopwatch',
  'download',
  'folder',
  'file',
  'hidden',
  'perf-replication',
  'role',
  'visible',
  'information-reversed',
  'true',
  'false',
  'upload',
  'control-lock',
  'edition-enterprise',
  'edition-oss',
  'check-plain',
  'check-circle-fill',
  'cancel-square-outline',
  'status-indicator',
  'person',
  'console',
  'checkmark-circled-outline',
  'close-circled-outline',
  'neutral-circled-outline',
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
