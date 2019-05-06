/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './popup-menu.md';


storiesOf('PopupMenu/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`PopupMenu`, () => ({
    template: hbs`
        <h5 class="title is-5">Popup Menu</h5>
        <PopupMenu/>
    `,
    context: {},
  }),
  {notes}
);
