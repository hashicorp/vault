import { attribute } from 'ember-cli-page-object';

export default {
  viewBox: attribute('viewBox', '[data-test-radial-progress]'),
  height: attribute('height', '[data-test-radial-progress]'),
  width: attribute('width', '[data-test-radial-progress]'),
  cx: attribute('cx', '[data-test-path]'),
  cy: attribute('cy', '[data-test-path]'),
  r: attribute('r', '[data-test-path]'),
  strokeWidth: attribute('stroke-width', '[data-test-path]'),
  strokeDash: attribute('stroke-dasharray', '[data-test-progress]'),
  strokeDashOffset: attribute('stroke-dashoffset', '[data-test-progress]'),
};
