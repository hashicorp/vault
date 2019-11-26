import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';
import notes from './form-save-buttons.md';

storiesOf('FormSaveButtons/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add(
    `FormSaveButtons`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Form save buttons</h5>
      <FormSaveButtons
        @isSaving={{this.save}}
        @saveButtonText={{this.saveButtonText}}
        @cancelButtonText={{this.cancelButtonText}}
        @includeBox={{this.includeBox}}
        @onCancel={{this.onCancel}}
        />
    `,

      context: {
        save: boolean('saving?', false),
        includeBox: boolean('include box?', true),
        saveButtonText: text('save button text', 'Save'),
        cancelButtonText: text('cancel button text', 'Cancel'),
        onCancel: () => {
          console.log('Canceled!');
        },
      },
    }),
    { notes }
  );
