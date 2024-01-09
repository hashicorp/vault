import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';

export class UserState {
  @tracked firstName = '';
  @tracked lastName = '';

  @tracked emailAddress = '';

  @tracked subscribed = false;

  @tracked department = null;
  @tracked manager = '';

  // constructor(firstName, lastName, emailAddress, subscribed, department, manager) {
  // this.firstName = firstName || this.firstName;
  // this.lastName = lastName || this.lastName;
  // this.emailAddress = emailAddress || this.emailAddress;
  // this.subscribed = subscribed || this.subscribed;
  // this.department = department || this.department;
  // this.manager = manager || this.manager;
  // }

  get fullName() {
    if (!this.firstName || !this.lastName) return '';
    return this.firstName + ' ' + this.lastName;
  }
}

const failPromise = () =>
  new Promise((resolve, reject) => {
    reject('Something else went wrong');
  });

export default class ExampleFormComponent extends Component {
  @tracked userState = new UserState();
  @tracked errors;
  @tracked successful = false;

  get allByKey() {
    return {
      firstName: {
        name: 'firstName',
      },
      lastName: {
        name: 'lastName',
      },
      department: { name: 'department' },
    };
  }
  get fields() {
    if (this.userState.subscribed === 'Yes') {
      return ['firstName', 'lastName', 'department'];
    } else {
      return ['firstName', 'lastName'];
    }
  }

  get mainFields() {
    return [
      {
        name: 'firstName',
        label: 'First Name',
        helpText: 'Your first name',
        editType: 'text',
        required: true,
        placeholder: 'John',
      },
      {
        name: 'lastName',
        label: 'Last Name',
        helpText: 'Your last name',
        editType: 'text',
        required: true,
        placeholder: 'Doe',
      },
      {
        name: 'emailAddress',
        label: 'Email Address',
        helpText: 'Your email address',
        editType: 'text',
        type: 'email',
        required: true,
        placeholder: '',
      },
    ];
  }

  get settingFields() {
    return [
      {
        name: 'subscribed',
        label: 'Subscribed',
        editType: 'select',
        required: true,
        options: ['Yes', 'No'],
      },
      {
        name: 'department',
        label: 'Department',
        helpText: 'choose your department',
        editType: 'select',
        required: true,
        options: ['Engineering', 'Marketing', 'Sales', 'Support'],
      },
      {
        name: 'manager',
        label: 'Manager',
        helpText: 'Select your manager',
        editType: 'select',
        required: true,
        options: ['John Doe', 'Jane Doe', 'John Smith', 'Jane Smith'],
      },
    ];
  }

  @action formChange(key, value) {
    this.userState[key] = value;
  }

  @task *submitForm() {
    this.errors = null;
    this.successful = false;
    yield timeout(1000);
    return failPromise();
  }
}
