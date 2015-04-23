Demo.Step = DS.Model.extend({
  name: DS.attr('string'),
  humanName: DS.attr('string'),

  instructionTemplate: Ember.computed.alias('name')
});
