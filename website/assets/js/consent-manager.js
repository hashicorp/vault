import { open, init } from '@hashicorp/hashi-consent-manager'

window.openConsentManager = () => open()

init({
  version: 1,
  container: '#consent-manager',
  companyName: 'HashiCorp',
  privacyPolicyLink: '/privacy',
  segmentWriteKey: segmentWriteKey,
  utilServerRoot: utilityServerRoot,
  segmentServices: [
    {
      key: 'googleanalytics',
      name: 'Google Analytics',
      description:
        'Google Analytics is a popular service for tracking web traffic. We use this data to determine what content our users find important so that we can dedicate more resources toward it.',
      category: 'Analytics'
    },
    {
      name: 'Marketo V2',
      description:
        'Marketo is a marketing automation tool that allows us to segment users into different categories based off of their behaviors.  We use this information to provide tailored information to users in our email campaigns.',
      category: 'Email Marketing'
    },
    {
      name: 'Hull',
      description:
        'Hull is a tool that we use to clean up analytics data and send it between different services. It does not add any javascript tracking code to this site.',
      category: 'Analytics'
    },
    {
      name: 'Hotjar',
      description:
        'Hotjar is a service that generates heatmaps of where users click on our sites. We use this information to ensure that our site is not confusing, and simple to use and navigate.',
      category: 'Analytics'
    }
  ],
  categories: [
    {
      name: 'Functional',
      description:
        'Functional services provide a utility to the website, such as the ability to log in, or to get live support. Disabling any of these scripts will cause that utility to be missing from the site.'
    },
    {
      name: 'Analytics',
      description:
        'Analytics services keep track of page traffic and user behavior while browsing the site. We use this data internally to improve the usability and performance of the site. Disabling any of these scripts makes it more difficult for us to understand how our site is being used, and slower to improve it.'
    },
    {
      name: 'Email Marketing',
      description:
        'Email Marketing services track user behavior while browsing the site. We use this data internally in our marketing efforts to provide users contextually relevant information based off of their behaviors. Disabling any of these scripts makes it more difficult for us to provide you contextually relevant information.'
    }
  ],
  additionalServices: [
    {
      name: 'OptinMonster',
      description:
        "OptinMonster is a service that we use to show a prompt to sign up for our newsletter if it's perceived that you are interested in our content.",
      category: 'Functional',
      body: `var om597a24292a958,om597a24292a958_poll=function(){var e=0;return function(t,a){clearInterval(e),e=setInterval(t,a)}}();!function(e,t,a){if(e.getElementById(a))om597a24292a958_poll(function(){if(window.om_loaded&&!om597a24292a958)return(om597a24292a958=new OptinMonsterApp).init({s:"35109.597a24292a958",staging:0,dev:0,beta:0})},25);else{var n=!1,o=e.createElement("script");o.id=a,o.src="//a.optnmstr.com/app/js/api.min.js",o.async=!0,o.onload=o.onreadystatechange=function(){if(!(n||this.readyState&&"loaded"!==this.readyState&&"complete"!==this.readyState))try{n=om_loaded=!0,(om597a24292a958=new OptinMonsterApp).init({s:"35109.597a24292a958",staging:0,dev:0,beta:0}),o.onload=o.onreadystatechange=null}catch(e){}},(document.getElementsByTagName("head")[0]||document.documentElement).appendChild(o)}}(document,0,"omapi-script");analytics.page()`
    }
  ]
})
