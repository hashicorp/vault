/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object, text, boolean } from '@storybook/addon-knobs';
import notes from './select.md';

const OPTIONS = [
  { value: 'mon', label: 'Monday', spanish: 'lunes' },
  { value: 'tues', label: 'Tuesday', spanish: 'martes' },
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
          @selectedValue={{selectedValue}}
        />
    `,
      context: {
        options: object('options', OPTIONS),
        label: text('label', 'Favorite fruit'),
        isFullwidth: boolean('isFullwidth', false),
        isInline: boolean('isInline', false),
        valueAttribute: text('valueAttribute', 'value'),
        labelAttribute: text('labelAttribute', 'label'),
        selectedValue: text('selectedValue', OPTIONS[1].value),
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
            @isInline={{true}}/>
        </Toolbar>
    `,
      context: {
        label: text('label', 'Favorite fruit'),
        options: object('options', OPTIONS),
        valueAttribute: text('valueAttribute', 'value'),
        labelAttribute: text('labelAttribute', 'label'),
      },
    }),
    { notes }
  );
