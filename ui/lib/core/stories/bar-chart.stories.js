import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs } from '@storybook/addon-knobs';

storiesOf('BarChart', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`BarChart`, () => ({
    template: hbs`
      <h5 class="title is-5">Bar Chart</h5>
      <BarChart/>
    `,
    context: {},
  }));
