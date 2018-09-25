import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  'data-test-radial-progress': true,
  tagName: 'svg',
  classNames: 'radial-progress',
  attributeBindings: ['size:width', 'size:height', 'viewBox'],
  progressDecimal: null,
  size: 20,
  strokeWidth: 1,

  viewBox: computed('size', function() {
    let s = this.get('size');
    return `0 0 ${s} ${s}`;
  }),
  centerValue: computed('size', function() {
    return this.get('size') / 2;
  }),
  r: computed('size', 'strokeWidth', function() {
    return (this.get('size') - this.get('strokeWidth')) / 2;
  }),
  c: computed('r', function() {
    return 2 * Math.PI * this.get('r');
  }),
  dashArrayOffset: computed('c', 'progressDecimal', function() {
    return this.get('c') * (1 - this.get('progressDecimal'));
  }),
});
