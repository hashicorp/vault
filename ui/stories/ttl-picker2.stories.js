import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './ttl-picker2.md';
import { withKnobs, text, boolean, select } from '@storybook/addon-knobs';

storiesOf('TtlPicker2/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `TtlPicker2`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Ttl Picker2</h5>
      <TtlPicker2 @unit="h" />
    `,
      context: {
        enabled: boolean('enabled', false),
      },
    }),
    { notes }
  );
