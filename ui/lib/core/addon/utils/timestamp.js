import Ember from 'ember';

/**
 * Use this instead of new Date() throughout the app so that time-related logic is easier to test.
 */
const timestamp = {
  // Method defined within an object so it can be stubbed
  /**
   * * Use timestamp.now to create a date for the current moment. In testing context, it always returns Apr 3, 2018 at 14:15:30
   * @returns Date object
   */
  now: () => {
    if (Ember.testing) {
      // April 3, 2018 -- when we moved the Vault UI assets to OSS
      return new Date('2018-04-03T14:15:30');
    }
    return new Date();
  },
};

export default timestamp;
