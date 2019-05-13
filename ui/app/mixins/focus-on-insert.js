import { schedule } from '@ember/runloop';
import { on } from '@ember/object/evented';
import Mixin from '@ember/object/mixin';

export default Mixin.create({
  // selector passed to `this.$()` to find the element to focus
  // defaults to `'input'`
  focusOnInsertSelector: null,
  shouldFocus: true,

  // uses Ember.on so that we don't have to worry about calling _super if
  // didInsertElement is overridden
  focusOnInsert: on('didInsertElement', function() {
    schedule('afterRender', this, 'focusOnInsertFocus');
  }),

  focusOnInsertFocus() {
    if (this.get('shouldFocus') === false) {
      return;
    }
    this.forceFocus();
  },

  forceFocus() {
    var $targ = this.element.querySelectorAll(this.get('focusOnInsertSelector') || 'input[type="text"]')[0];
    if ($targ && $targ !== document.activeElement) {
      $targ.focus();
    }
  },
});
