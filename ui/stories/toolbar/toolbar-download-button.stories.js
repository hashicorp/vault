/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text } from '@storybook/addon-knobs';
import notes from './toolbar-download-button.md';


storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ToolbarDownloadButton`,() => ({
    template: hbs`
      <h5 class="title is-5">ToolbarLink</h5>
      <div style="width: 400px;">
        <Toolbar>
          <ToolbarActions>
            <ToolbarDownloadButton
              @actionText={{label}}
            />
          </ToolbarActions>
        </Toolbar>
      </div>
    `,
    context: {
      label: text('Button text', 'Download policy'),
    },
  }),
  {notes}
);
