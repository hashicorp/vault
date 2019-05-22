/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './http-requests-table.md';


storiesOf('HttpRequestsTable/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`HttpRequestsTable`, () => ({
    template: hbs`
        <h5 class="title is-5">Http Requests Table</h5>
        <HttpRequestsTable @counters={{counters}}/>
    `,
    context: {
      counters: [
        {
          start_time: '2019-05-01T00:00:00Z',
          total: 50,
        },
        {
          start_time: '2019-04-01T00:00:00Z',
          total: 45,
        },
        {
          start_time: '2019-03-01T00:00:00Z',
          total: 55,
        },
      ]
    },
  }),
  { notes }
);
