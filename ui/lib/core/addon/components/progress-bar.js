import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/progress-bar';

/**
 * @module ProgressBar
 * `ProgressBar` components are a rectangular bar showing the progress of something.
 *
 * @example
 * ```js
 * <ProgressBar @threshold={{}} @progress={{}}/>
 * ```
 *
 * @param type=null {String} - The banner type. This comes from the message-types helper.
 * @param [message=null {String}] - The message to display within the banner.
 *
 */

export default Component.extend({
  layout,
  tagName: '',
  threshold: null,
  progress: null,
  progressPercent: computed('threshold', 'progress', function() {
    const { threshold, progress } = this;
    if (threshold && progress) {
      const stuff = (progress / threshold) * 100;
      console.log(stuff);
      return (progress / threshold) * 100;
    }
    return 0;
  }),
});
