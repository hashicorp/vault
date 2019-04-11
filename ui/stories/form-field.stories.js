/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './form-field.md';
import EmberObject from '@ember/object';

const createAttr = (name, type, options) => {
  return {
    name,
    type,
    options,
  };
};

storiesOf('Form/FormField/', module)
  .add(
    `FormField|string`,
    () => ({
      template: hbs`
        <h5 class="title is-5">String Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'string', { defaultValue: 'default' }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|boolean`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Boolean Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'boolean', { defaultValue: false }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|number`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Number Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'number', { defaultValue: 5 }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|object`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Object Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'object'),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|textarea`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Textarea Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'string', { defaultValue: 'goodbye', editType: 'textarea' }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|file`,
    () => ({
      template: hbs`
        <h5 class="title is-5">File Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'string', { editType: 'file' }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|ttl`,
    () => ({
      template: hbs`
        <h5 class="title is-5">ttl Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', null, { editType: 'ttl' }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|stringArray`,
    () => ({
      template: hbs`
        <h5 class="title is-5">String Array Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('foo', 'string', { editType: 'stringArray' }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  )
  .add(
    `FormField|sensitive`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Sensitive Form Field</h5>
        <FormField @attr={{attr}} @model={{model}}/>
    `,
      context: {
        attr: createAttr('password', 'string', { sensitive: true }),
        model: EmberObject.create({}),
      },
    }),
    { notes }
  );
