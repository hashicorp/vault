/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-popup.md';
import { messageTypes } from '../app/helpers/message-types.js';

storiesOf('AlertPopup/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    `AlertPopup`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Alert Popup</h5>
        <AlertPopup
          @type={{info}}
          @message="Hello!"
          @close={{close}}/>
          <h5 class="title is-5">Alert Popup</h5>
        <AlertPopup
          @type={{success}}
          @message="Hello!"
          @close={{close}}/>
          <h5 class="title is-5">Alert Popup</h5>
        <AlertPopup
          @type={{danger}}
          @message="Hello!"
          @close={{close}}/>
          <h5 class="title is-5">Alert Popup</h5>
        <AlertPopup
          @type={{warning}}
          @message="Hello!"
          @close={{close}}/>
    `,
      context: {
        close: () => {
          console.log('closing!');
        },
        info: messageTypes(['info']),
        success: messageTypes(['success']),
        danger: messageTypes(['danger']),
        warning: messageTypes(['warning']),
      },
    }),
    { notes }
  );
