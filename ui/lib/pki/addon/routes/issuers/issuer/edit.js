import { withConfirmLeave } from 'core/decorators/confirm-leave';
import PkiIssuerDetailsRoute from './details';

withConfirmLeave();
export default class PkiIssuerEditRoute extends PkiIssuerDetailsRoute {}
