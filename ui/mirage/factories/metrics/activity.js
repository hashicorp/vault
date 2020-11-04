import Mirage from 'ember-cli-mirage';

export default Mirage.Factory.extend({
  endTime: '2020-09-30T23:59:59Z',
  end_time: '2020-09-30T23:59:59Z',
  startTime: '2020-03-01T00:00:00Z',
  start_time: '2020-03-01T00:00:00Z',
  total: {
    clients: 900,
    distinct_entities: 462,
    non_entity_tokens: 461,
  },
  request_id(i) {
    return `00000${i}`;
  },
});
