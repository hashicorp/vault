/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './http-requests-table.md';


storiesOf('HttpRequests/Table/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`HttpRequestsTable`, () => ({
    template: hbs`
        <h5 class="title is-5">Http Requests Table</h5>
        <HttpRequestsTable @counters={{counters}}/>
    `,
    context: {
      counters: [
        {
          start_time: '2019-04-01T05:00:00.000Z',
          total: 5500,
        },
        {
          start_time: '2019-05-01T05:00:00.000Z',
          total: 4500,
        },
        {
          start_time: '2019-06-01T05:00:00.000Z',
          total: 5000,
        },
      ]
    },
  }),
  { notes }
);
