/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object, text, boolean, select } from '@storybook/addon-knobs';
import notes from './select.md';

const OPTIONS = [
  { value: 'mon', label: 'Monday', spanish: 'lunes' },
  { value: 'tues', label: 'Tuesday', spanish: 'martes' },
  { value: 'weds', label: 'Wednesday', spanish: 'miercoles' },
];

storiesOf('Select/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `Select`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Select</h5>
        <Select
          @options={{options}}
          @label={{label}}
          @isInline={{isInline}}
          @isFullwidth={{isFullwidth}}
          @valueAttribute={{valueAttribute}}
          @labelAttribute={{labelAttribute}}
          @selectedValue={{selectedValue}}
        />
    `,
      context: {
        options: object('options', OPTIONS),
        label: text('label', 'Day of the week'),
        isFullwidth: boolean('isFullwidth', false),
        isInline: boolean('isInline', false),
        valueAttribute: select('valueAttribute', Object.keys(OPTIONS[0]), 'value'),
        labelAttribute: select('labelAttribute', Object.keys(OPTIONS[0]), 'label'),
        selectedValue: select('selectedValue', OPTIONS.map(o => o.label), 'tues'),
      },
    }),
    { notes }
  )
  .add(
    `Select in a Toolbar`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Select</h5>
        <Toolbar>
          <Select
            @options={{options}}
            @label={{label}}
            @valueAttribute={{valueAttribute}}
            @labelAttribute={{labelAttribute}}
            @isInline={{true}}/>
        </Toolbar>
    `,
      context: {
        label: text('label', 'Day of the week'),
        options: object('options', OPTIONS),
        valueAttribute: select('valueAttribute', Object.keys(OPTIONS[0]), 'value'),
        labelAttribute: select('labelAttribute', Object.keys(OPTIONS[0]), 'label'),
      },
    }),
    { notes }
  );
