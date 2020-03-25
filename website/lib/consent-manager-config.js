const isProd = process.env.NODE_ENV === 'production'

const segmentWriteKey = isProd
  ? 'OdSFDq9PfujQpmkZf03dFpcUlywme4sC'
  : '0EXTgkNx0Ydje2PGXVbRhpKKoe5wtzcE'

// TODO: refactor into web components
let utilityServerRoot = isProd
  ? 'https://util.hashicorp.com'
  : 'https://hashicorp-web-util-staging.herokuapp.com'

if (process.env.UTIL_SERVER) {
  utilityServerRoot = process.env.UTIL_SERVER.replace(/\/$/, '')
}

// Consent manager configuration
export default {
  version: 3,
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
      name: 'Hotjar',
      description:
        'Hotjar is a service that generates heatmaps of where users click on our sites. We use this information to ensure that our site is not confusing, and simple to use and navigate.',
      category: 'Analytics'
    },
    {
      name: 'LinkedIn Insight Tag',
      description:
        'This small script allows us to see how effective our linkedin campaigns are by showing which users have clicked through to our site.',
      category: 'Analytics'
    },
    {
      name: 'Marketo V2',
      description:
        'Marketo is a marketing automation tool that allows us to segment users into different categories based off of their behaviors.  We use this information to provide tailored information to users in our email campaigns.'
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
      body: `var om598c8e3a6e43d,om598c8e3a6e43d_poll=function(){var r=0;return function(n,l){clearInterval(r),r=setInterval(n,l)}}();!function(e,t,n){if(e.getElementById(n)){om598c8e3a6e43d_poll(function(){if(window['om_loaded']){if(!om598c8e3a6e43d){om598c8e3a6e43d=new OptinMonsterApp();return om598c8e3a6e43d.init({"s":"35109.598c8e3a6e43d","staging":0,"dev":0,"beta":0});}}},25);return;}var d=false,o=e.createElement(t);o.id=n,o.src="https://a.optnmstr.com/app/js/api.min.js",o.async=true,o.onload=o.onreadystatechange=function(){if(!d){if(!this.readyState||this.readyState==="loaded"||this.readyState==="complete"){try{d=om_loaded=true;om598c8e3a6e43d=new OptinMonsterApp();om598c8e3a6e43d.init({"s":"35109.598c8e3a6e43d","staging":0,"dev":0,"beta":0});o.onload=o.onreadystatechange=null;}catch(t){}}}};(document.getElementsByTagName("head")[0]||document.documentElement).appendChild(o)}(document,"script","omapi-script");`
    }
  ]
}
