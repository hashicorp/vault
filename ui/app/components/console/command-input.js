import Ember from 'ember';
import keys from 'vault/lib/keycodes';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';


export default Ember.Component.extend({
  onExecuteCommand() {},
  onValueUpdate() {},
  value: null,

  actions: {
    handleKeyUp: function(event) {
      var keyCode = event.keyCode;
      if (keyCode === keys.ENTER) {
        this.get('onExecuteCommand')(event.target.value);
      } else {
        this.get('onValueUpdate')(event.target.value);
      }
    }
  }
});

// up cycles through history
// down clears if there is no more history or cycles down through history
//
