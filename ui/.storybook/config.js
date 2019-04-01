import { configure, addParameters, addDecorator } from '@storybook/ember';
import { INITIAL_VIEWPORTS } from '@storybook/addon-viewport';
import theme from './theme.js';

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
  Object.assign(element.style, styles.style);

  const innerElement = document.createElement('div');

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
