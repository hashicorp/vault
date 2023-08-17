import { ConsentManagerService } from '@hashicorp/react-consent-manager/types'

const localConsentManagerServices: ConsentManagerService[] = [
  {
    name: 'Qualified Chatbot',
    description:
      'Qualified is a chatbot service that allows visitors to chat with our sales staff through the website.',
    category: 'Email Marketing',
    url: 'https://js.qualified.com/qualified.js?token=CWQA3q9CaEKHNF2t',
    async: true,
  },
]

export default localConsentManagerServices
