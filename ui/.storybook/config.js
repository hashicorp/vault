import { configure, addParameters, addDecorator } from '@storybook/ember';
import { INITIAL_VIEWPORTS } from '@storybook/addon-viewport';
import theme from './theme.js';

function loadStories() {
  // automatically import all files ending in *.stories.js
  const appStories = require.context('../stories', true, /.stories.js$/);
  const addonAndRepoStories = require.context('../lib', true, /.stories.js$/);
  appStories.keys().forEach(filename => appStories(filename));
  addonAndRepoStories.keys().forEach(filename => addonAndRepoStories(filename));
}

addParameters({
  viewport: { viewports: INITIAL_VIEWPORTS },
  options: { theme },
});

addDecorator(storyFn => {
  const { template, context } = storyFn();

  // This adds styling to the Canvas tab.
  const styles = {
    style: {
      margin: '20px',
    },
  };

  // Create a div to wrap the Canvas tab with the applied styles.
  const element = document.createElement('div');
  Object.assign(element.style, styles.style);

  const innerElement = document.createElement('div');
  const wormhole = document.createElement('div');
  wormhole.setAttribute('id', 'ember-basic-dropdown-wormhole');
  innerElement.appendChild(wormhole);

  element.appendChild(innerElement);
  innerElement.appendTo = function appendTo(el) {
    el.appendChild(element);
  };

  return {
    template,
    context,
    element: innerElement,
  };
});

configure(loadStories, module);
