/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object } from '@storybook/addon-knobs';
import notes from './http-requests-container.md';

const COUNTERS = [
  { start_time: '2017-04-01T00:00:00Z', total: 5500 },
  { start_time: '2018-06-01T00:00:00Z', total: 4500 },
  { start_time: '2018-07-01T00:00:00Z', total: 4500 },
  { start_time: '2018-08-01T00:00:00Z', total: 6500 },
  { start_time: '2018-09-01T00:00:00Z', total: 5500 },
  { start_time: '2018-10-01T00:00:00Z', total: 4500 },
  { start_time: '2018-11-01T00:00:00Z', total: 6500 },
  { start_time: '2018-12-01T00:00:00Z', total: 5500 },
  { start_time: '2019-01-01T00:00:00Z', total: 2500 },
  { start_time: '2019-02-01T00:00:00Z', total: 3500 },
  { start_time: '2019-03-01T00:00:00Z', total: 5000 },
];

storiesOf('HttpRequests/Container/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `HttpRequestsContainer`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Http Requests Container</h5>
        <HttpRequestsContainer @counters={{counters}}/>
    `,
      context: {
        counters: object('counters', COUNTERS),
      },
    }),
    { notes }
  );
