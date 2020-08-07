import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, boolean, object } from '@storybook/addon-knobs';
import notes from './selectable-card-container.md';

const MODEL = {
  totalEntities: 0,
  httpsRequests: [
    { start_time: '2018-12-01T00:00:00Z', total: 5500 },
    { start_time: '2019-01-01T00:00:00Z', total: 4500 },
    { start_time: '2019-02-01T00:00:00Z', total: 5000 },
    { start_time: '2019-03-01T00:00:00Z', total: 5000 },
  ],
  totalTokens: 1,
};
const GRID_CONTAINER = true;

storiesOf('SelectableCard/SelectableCardContainer', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `SelectableCardContainer`,
    () => ({
      template: hbs`
      <SelectableCardContainer @counters={{model}} @gridContainer={{gridContainer}} />
        
    `,
      context: {
        model: object('model', MODEL),
        gridContainer: boolean('gridContainer', GRID_CONTAINER),
      },
    }),
    { notes }
  );
