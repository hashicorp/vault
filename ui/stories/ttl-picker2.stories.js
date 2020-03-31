import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './ttl-picker2.md';
import { withKnobs, text, boolean, select } from '@storybook/addon-knobs';

storiesOf('TtlPicker2/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `TtlPicker2|single`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Ttl Picker2</h5>
      <TtlPicker2
        @unit={{unit}}
        @enableTTL={{enableTTL}}
        @label={{label}}
        @helperTextDisabled={{helperTextDisabled}}
        @helperTextEnabled={{helperTextEnabled}}
        @onChange={{onChange}}
        @time={{time}}
      />
    `,
      context: {
        enableTTL: boolean('enableTTL', false),
        unit: select('unit', ['s', 'm', 'h', 'd'], 'm'),
        time: text('time', '40'),
        label: text('label', 'Main label of TTL'),
        helperTextDisabled: text('helperTextDisabled', 'This helper text displays when TTL is disabled'),
        helperTextEnabled: text('helperTextEnabled', 'Enabling TTL will show this text instead'),
        onChange: ttl => {
          console.log('onChange fired', ttl);
        },
      },
    }),
    { notes }
  )
  .add(
    `TtlPicker2|nested`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Ttl Picker2</h5>
      <TtlPicker2 @unit="m" @enableTTL={{enableTTL}}>
        <TtlPicker2
          @label="Maximum time to live (Max TTL)"
          @helperTextDisabled="Allow tokens to be renewed indefinitely"
          @unit="h"
        />
      </TtlPicker2>
    `,
      context: {
        enableTTL: boolean('enableTTL', true),
      },
    }),
    { notes }
  );
