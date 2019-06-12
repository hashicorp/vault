/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object } from '@storybook/addon-knobs';
import notes from './http-requests-table.md';

const COUNTERS = [
  { start_time: '2019-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

storiesOf('HttpRequests/Table/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs()
  )
  .add(`HttpRequestsTable`, () => ({
    template: hbs`
        <h5 class="title is-5">Http Requests Table</h5>
        <HttpRequestsTable @counters={{counters}}/>
    `,
    context: {
      counters: object('counters', COUNTERS),
    }
  }),
  { notes }
);
