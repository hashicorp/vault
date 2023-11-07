import Component from '@glimmer/component';
import { withAuthForm } from 'vault/decorators/auth-form';

@withAuthForm('userpass')
export default class UserpassComponent extends Component {}
