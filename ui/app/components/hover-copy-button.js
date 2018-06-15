import Ember from 'ember';

export default Ember.Component.extend({
  classNameBindings: 'alwaysShow:hover-copy-button-static:hover-copy-button',
  copyValue: null,
  alwaysShow: false,

  tooltipText: 'Copy',
});
