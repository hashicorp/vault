import { configure, addParameters, addDecorator } from '@storybook/ember';
import { INITIAL_VIEWPORTS } from '@storybook/addon-viewport';
import theme from './theme.js';
import { assign } from '@ember/polyfills';

function loadStories() {
  // automatically import all files ending in *.stories.js
  const req = require.context('../stories/', true, /.stories.js$/);
  req.keys().forEach(filename => req(filename));
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
  assign(element.style, styles.style);

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
