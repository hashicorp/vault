/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './chevron.md';


storiesOf('Chevron/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`Chevron`, () => ({
    template: hbs`
        <h5 class="title is-5">Chevron</h5>
        <Chevron @direction={{direction}} />
    `,
    context: {
      direction: select('Direction', ['right', 'down', 'left', 'up'], 'right'),
    },
  }),
  {notes}
);
