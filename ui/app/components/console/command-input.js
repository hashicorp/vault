import Ember from 'ember';
import keys from 'vault/lib/keycodes';

export default Ember.Component.extend({
  'data-test-component': 'console/command-input',
  classNames: 'console-ui-input',
  onExecuteCommand() {},
  onValueUpdate() {},
  onShiftCommand() {},
  value: null,

  didRender() {
    this.element.scrollIntoView();
  },
  actions: {
    handleKeyUp: function(event) {
      var keyCode = event.keyCode;
      switch (keyCode) {
        case keys.ENTER:
          this.get('onExecuteCommand')(event.target.value);
          break;
        case keys.UP:
        case keys.DOWN:
          this.get('onShiftCommand')(keyCode);
          break;
        default:
          this.get('onValueUpdate')(event.target.value);
      }
    },
  },
});

// up cycles through history
// down clears if there is no more history or cycles down through history
//
