import Mirage from 'ember-cli-mirage';

export default Mirage.Factory.extend({
  end_time: '2020-10-31T23:59:59Z',
  start_time: '2020-09-01T00:00:00Z',
  total: function() {
    return {
      clients: 20,
      distinct_entities: 10,
      non_entity_tokens: 10,
    };
  },
});
