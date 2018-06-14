import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

const { computed } = Ember;
const GLYPHS_WITH_SVG_TAG = [
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
];

export default Ember.Component.extend({
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
    return GLYPHS_WITH_SVG_TAG.includes(this.get('glyph'));
  }),

  size: computed(function() {
    return 12;
  }),

  partialName: computed('glyph', function() {
    const glyph = this.get('glyph');
    return `svg/icons/${Ember.String.camelize(glyph)}`;
  }),
});
