/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, boolean, text } from '@storybook/addon-knobs';
import notes from './info-table-row.md';


storiesOf('InfoTableRow/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({escapeHTML: false}))
  .add(`InfoTableRow with text value`, () => ({
    template: hbs`
      <h5 class="title is-5">Info Table Row</h5>
      <InfoTableRow @value={{value}} @label={{label}} @alwaysRender={{alwaysRender}} />
    `,
    context: {
      label: text('Label', 'TTL'),
      value: text('Value', '30m'),
      alwaysRender: boolean('Always render?', false),
    },
  }),
  {notes}
  )
  .add(`InfoTableRow with boolean value`, () => ({
    template: hbs`
      <h5 class="title is-5">Info Table Row</h5>
      <InfoTableRow @value={{value}} @label={{label}} @alwaysRender={{alwaysRender}} />
    `,
    context: {
      label: 'Local mount?',
      value: boolean('Value', true),
      alwaysRender: boolean('Always render?', true),
    },
  }),
  {notes}
);
