import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { text, withKnobs } from '@storybook/addon-knobs';
import notes from './stat-text.md';

storiesOf('StatText', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `StatText`,
    () => ({
      template: hbs`
    <h5 class="title is-5">StatText Component</h5>
    <StatText
     @label={{label}}
     @value={{value}}
     @size={{size}}
     @subText={{subText}} />
    `,
      context: {
        label: text('label', 'Active Clients'),
        value: text('value', '4,198'),
        size: text('size', 'l'),
        subText: text('subText', 'These are your active clients'),
      },
    }),
    { notes }
  );
