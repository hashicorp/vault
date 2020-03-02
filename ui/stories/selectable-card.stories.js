import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, boolean, number, text } from '@storybook/addon-knobs';
import notes from './selectable-card.md';

const CARD_TITLE = 'Tokens';
const SUB_TEXT = 'Total';
const TOTAL_HTTP_REQUESTS = 100;

storiesOf('SelectableCard/SelectableCard', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `SelectableCard`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Selectable Card</h5>
        <SelectableCard @cardTitle={{cardTitle}} @total={{totalHttpRequests}} @subText={{subText}} />
    `,
      context: {
        cardTitle: text('cardTitle', CARD_TITLE),
        totalHttpRequests: number('totalHttpRequests', TOTAL_HTTP_REQUESTS),
        subText: text('subText', SUB_TEXT),
      },
    }),
    { notes }
  );
