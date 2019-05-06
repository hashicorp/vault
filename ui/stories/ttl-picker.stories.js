/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './ttl-picker.md';


storiesOf('TtlPicker/', module)
  .addParameters({ options: { showPanel: false } })
  .add(`TtlPicker`, () => ({
    template: hbs`
      <h5 class="title is-5">Ttl Picker</h5>
      <TtlPicker />
    `,
  }),
  {notes}
);
