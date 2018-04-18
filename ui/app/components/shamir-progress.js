import Ember from 'ember';

export default Ember.Component.extend({
  threshold: null,
  progress: null,
  classNames: ['shamir-progress'],
  progressPercent: Ember.computed('threshold', 'progress', function() {
    const { threshold, progress } = this.getProperties('threshold', 'progress');
    if (threshold && progress) {
      return progress / threshold * 100;
    }
    return 0;
  }),
});
