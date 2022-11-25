import { AuthState, onAuthUIStateChange } from "@aws-amplify/ui-components";
import { AmplifyAuthenticator } from "@aws-amplify/ui-react";
import { Amplify, Auth } from "aws-amplify";
import React from "react";
import Router from "./components/Router";

interface User {
	username: string;
}

const App: React.FunctionComponent = () => {
	const awsconfig = {
		Auth: {
			region: process.env.REACT_APP_AWS_REGION,
			userPoolId: process.env.REACT_APP_USER_POOL_ID,
			userPoolWebClientId: process.env.REACT_APP_API_CLIENT_ID,
		},
		API: {
			endpoints: [
				{
					name: "photosAPIdev",
					endpoint: `https://${process.env.REACT_APP_BASE_DOMAIN}/api`,
					custom_header: async () => {
						return {
							Authorization: `Bearer ${(await Auth.currentSession())
								.getIdToken()
								.getJwtToken()}`,
						};
					},
				},
			],
		},
	};
	Amplify.configure(awsconfig);

	const [authState, setAuthState] = React.useState<AuthState>();
	const [user, setUser] = React.useState<User | undefined>();

	React.useEffect(() => {
		return onAuthUIStateChange((nextAuthState, authData) => {
			setAuthState(nextAuthState);
			setUser(authData as User);
		});
	}, []);

	const isDevMode = process.env.NODE_ENV === "development";
	const isLoggedIn = authState === AuthState.SignedIn && user;

	return isDevMode || isLoggedIn ? (
		<div>
			<Router isDevMode={isDevMode} />
		</div>
	) : (
		<div>
			<h1>{process.env.NODE_ENV}</h1>
			<AmplifyAuthenticator />
		</div>
	);
};

export default App;
