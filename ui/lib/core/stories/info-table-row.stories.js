import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, boolean, text } from '@storybook/addon-knobs';
import notes from './info-table-row.md';

storiesOf('InfoTable/InfoTableRow', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({ escapeHTML: false }))
  .add(
    `InfoTableRow with text value`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Info Table Row</h5>
      <InfoTableRow 
        @value={{value}} 
        @label={{label}} 
        @helperText={{helperText}} 
        @alwaysRender={{alwaysRender}} 
        @defaultShown={{defaultShown}} 
        @tooltipText={{tooltipText}}
        @isTooltipCopyable={{isTooltipCopyable}} />
    `,
      context: {
        label: text('Label', 'TTL'),
        value: text('Value', '30m (hover to see the tooltip!)'),
        helperText: text('helperText', 'This is helperText - for a short description'),
        alwaysRender: boolean('Always render?', false),
        defaultShown: text('Default value', 'Some default value'),
        tooltipText: text('tooltipText', 'This is tooltipText'),
        isTooltipCopyable: boolean('Allow tooltip to be copied', true),
      },
    }),
    { notes }
  )
  .add(
    `InfoTableRow with boolean value`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Info Table Row</h5>
      <InfoTableRow 
        @value={{value}} 
        @label={{label}} 
        @helperText={{helperText}} 
        @alwaysRender={{alwaysRender}} />
    `,
      context: {
        label: 'Local mount?',
        value: boolean('Value', true),
        helperText: text('helperText', 'This is helperText - for a short description'),
        alwaysRender: boolean('Always render?', true),
      },
    }),
    { notes }
  );
