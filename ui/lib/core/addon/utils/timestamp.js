/**
 * Use this instead of new Date() throughout the app so that time-related logic is easier to test.
 */
const timestamp = {
  // Method defined within an object so it can be stubbed
  /**
   * * Use timestamp.now to create a date for the current moment. In testing context, stub this method so it returns an expected value
   * @returns Date object
   */
  now: () => {
    return new Date();
  },
};

export default timestamp;
