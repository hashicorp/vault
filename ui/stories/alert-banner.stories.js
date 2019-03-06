/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, select, boolean } from '@storybook/addon-knobs';

const TYPE_OPTIONS = ['warning', 'info', 'success', 'danger'];

storiesOf('AlertBanner', module)
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add('type', () => ({
    template: hbs`
      <AlertBanner @type={{type}} @message={{message}}/>
    `,
    context: {
      type: select('type', TYPE_OPTIONS, 'warning'),
      message: text('message', 'Here is a message.'),
    },
  }));
