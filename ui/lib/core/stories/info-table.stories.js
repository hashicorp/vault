import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object, text } from '@storybook/addon-knobs';
import notes from './info-table.md';

const ITEMS = ['https://127.0.0.1:8201', 'hello', 3];

storiesOf('InfoTable/InfoTable', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({ escapeHTML: false }))
  .add(
    `InfoTable`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Info Table</h5>
      <InfoTable
        @header={{header}}
        @items={{items}}
      />
    `,
      context: {
        header: text('Header', 'Column Header'),
        items: object('Items', ITEMS),
      },
    }),
    { notes }
  );
