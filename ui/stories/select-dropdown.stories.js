/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, array, text } from '@storybook/addon-knobs';
import notes from './select-dropdown.md';

const OPTIONS = ['apple', 'blueberry', 'cherry'];

storiesOf('SelectDropdown/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `SelectDropdown`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Select Dropdown</h5>
        <SelectDropdown @options={{options}} @dropdownLabel={{dropdownLabel}}/>

        <br />
        <br />
        <h5 class="title is-5">Select Dropdown - Inline</h5>
        <SelectDropdown @options={{options}} @dropdownLabel={{dropdownLabel}} @isInline={{true}}/>
    `,
      context: {
        options: array('options', OPTIONS),
        dropdownLabel: text('dropdownLabel', 'Favorite fruit'),
      },
    }),
    { notes }
  );
