/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';

storiesOf('Welcome', module).add('to Storybook', () => ({
  template: hbs`
  <AlertBanner @type="warning" @message={{warning}} />
      `,
  context: {
    warning: "Oops, don't do that again!",
  },
}));

storiesOf('Button', module)
  .add('with text', () => ({
    template: hbs`<button {{action onClick}}>Hello Button</button>`,
    context: {
      onClick: action('clicked'),
    },
  }))
  .add('with some emoji', () => ({
    template: hbs`
        <button {{action onClick}}>
          <span role="img" aria-label="so cool">
            😀 😎 👍 💯
          </span>
        </button>
      `,
    context: {
      onClick: action('clicked'),
    },
  }));
