import React from "react";
import ReactDOM from "react-dom";
import Amplify,{Auth} from 'aws-amplify'
import App from "./App";

const awsconfig =  {
  Auth: {
    region: process.env.REACT_APP_AWS_REGION,
    userPoolId: process.env.REACT_APP_USER_POOL_ID,
    userPoolWebClientId: process.env.REACT_APP_API_CLIENT_ID
  },
  API: {
    endpoints: [
      {
        name: "photosAPIdev",
        endpoint: "https://" + process.env.REACT_APP_BASE_DOMAIN + "/api",
        custom_header: async () => {
          return { Authorization: `Bearer ${(await Auth.currentSession()).getIdToken().getJwtToken()}` }
        }
      }
    ]
  }
};

Amplify.configure(awsconfig);

ReactDOM.render(
  <App />,
  document.getElementById("root")
);
