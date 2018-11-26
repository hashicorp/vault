// app/components/power-select/before-options.js
import { run } from '@ember/runloop';
import PowerSelectBeforeOptions from 'ember-power-select/components/power-select/before-options';

// TODO: Delete when this is fixed:
// https://github.com/cibernox/ember-power-select/issues/1130
export default PowerSelectBeforeOptions.extend({
  focusInput() {
    this.input = document.querySelector(
      `.ember-power-select-search-input[aria-controls="${this.listboxId}"]`
    );
    if (this.input) {
      run.next(this.input, 'focus');
    }
  },
});
