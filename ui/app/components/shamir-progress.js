import Ember from 'ember';
const { computed } = Ember;

export default Ember.Component.extend({
  threshold: null,
  progress: null,
  classNames: ['shamir-progress'],
  progressDecimal: computed('threshold', 'progress', function() {
    const { threshold, progress } = this.getProperties('threshold', 'progress');
    if (threshold && progress) {
      return progress / threshold;
    }
    return 0;
  }),
});
