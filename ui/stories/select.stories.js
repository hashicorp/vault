/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object, text, boolean } from '@storybook/addon-knobs';
import notes from './select.md';

const OPTIONS = ['apple', 'blueberry', 'cherry'];

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
        />
    `,
      context: {
        options: object('options', OPTIONS),
        valueAttribute: text('valueAttribute', 'value'),
        labelAttribute: text('labelAttribute', 'label'),
        label: text('label', 'Favorite fruit'),
        isFullwidth: boolean('isFullwidth', false),
        isInline: boolean('isInline', false),
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
        options: object('options', OPTIONS),
        valueAttribute: text('valueAttribute', 'value'),
        labelAttribute: text('labelAttribute', 'label'),
        label: text('label', 'Favorite fruit'),
      },
    }),
    { notes }
  );
