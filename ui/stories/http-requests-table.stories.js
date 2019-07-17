/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object } from '@storybook/addon-knobs';
import notes from './http-requests-table.md';

const COUNTERS = [
  { start_time: '2018-12-01T00:00:00Z', total: 5500 },
  { start_time: '2019-01-01T00:00:00Z', total: 4500 },
  { start_time: '2019-02-01T00:00:00Z', total: 5000 },
  { start_time: '2019-03-01T00:00:00Z', total: 5000 },
];

storiesOf('HttpRequests/Table/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `HttpRequestsTable`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Http Requests Table</h5>
        <HttpRequestsTable @counters={{counters}}/>
    `,
      context: {
        counters: object('counters', COUNTERS),
      },
    }),
    { notes }
  );
