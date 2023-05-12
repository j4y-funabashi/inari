import React from "react";
import Router from "./components/Router";

interface User {
	username: string;
}

const App: React.FunctionComponent<React.PropsWithChildren<unknown>> = () => {
	const isDevMode = process.env.NODE_ENV === "development";
	return (
		<div>
			<Router isDevMode={isDevMode} />
		</div>
	);
};

export default App;
