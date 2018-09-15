import Component from '@ember/component';
import keys from 'vault/lib/keycodes';

export default Component.extend({
  onExecuteCommand() {},
  onFullscreen() {},
  onValueUpdate() {},
  onShiftCommand() {},
  value: null,
  isFullscreen: null,

  actions: {
    handleKeyUp(event) {
      const keyCode = event.keyCode;
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
    fullscreen() {
      this.get('onFullscreen')();
    },
  },
});
