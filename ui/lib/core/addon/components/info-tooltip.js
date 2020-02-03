import Component from '@ember/component';
import layout from '../templates/components/info-tooltip';

export default Component.extend({
  layout,
  'data-test-component': 'info-tooltip',
  tagName: 'span',
  classNames: ['is-inline-block'],
});
