import Ember from 'ember';
import keys from 'vault/lib/keycodes';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';


export default Ember.Component.extend({
  tagName: 'input',
  onExecuteCommand() {},
  onValueUpdate() {},
  value: null,

  actions: {
    handleKeyUp: function(val, event) {
      var keyCode = event.keyCode;
      if (keyCode === keys.ENTER) {
        this.get('onExecuteCommand')(val);
      }
      this.get('onValueUpdate')(val);
    }

  }
});

// up cycles through history
// down clears if there is no more history or cycles down through history
//
