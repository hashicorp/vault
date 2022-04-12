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
  {
    name: 'Demandbase Tag',
    description:
      'The Demandbase tag is a tracking service to identify website visitors and measure interest on our website.',
    category: 'Analytics',
    url: 'https://tag.demandbase.com/960ab0a0f20fb102.min.js',
    async: true,
  },
]

export default localConsentManagerServices
