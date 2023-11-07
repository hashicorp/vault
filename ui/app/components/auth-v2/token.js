import Component from '@glimmer/component';
import { withAuthForm } from 'vault/decorators/auth-form';

@withAuthForm('token')
export default class TokenComponent extends Component {}
