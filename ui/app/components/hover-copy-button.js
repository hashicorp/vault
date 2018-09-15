import Component from '@ember/component';

export default Component.extend({
  'data-test-hover-copy': true,
  classNameBindings: 'alwaysShow:hover-copy-button-static:hover-copy-button',
  copyValue: null,
  alwaysShow: false,

  tooltipText: 'Copy',
});
