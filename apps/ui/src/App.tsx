import {AuthState, onAuthUIStateChange} from '@aws-amplify/ui-components';
import {AmplifyAuthenticator, AmplifySignOut} from '@aws-amplify/ui-react';
import React from 'react';
import MediaTimelinePage from './components/MediaTimelinePage'

interface User {
  username: string
}

const App: React.FunctionComponent = () => {
  const [authState, setAuthState] = React.useState<AuthState>();
  const [user, setUser] = React.useState<User | undefined>();

  React.useEffect(() => {
    return onAuthUIStateChange((nextAuthState, authData) => {
      setAuthState(nextAuthState);
      setUser(authData as User)
    });
  }, []);

  return authState === AuthState.SignedIn && user ? (
    <div>
      <nav>
        <div>Hello, {user.username}</div>
        <AmplifySignOut />
      </nav>
      <MediaTimelinePage />
    </div>
  ) : (
    <AmplifyAuthenticator />
  );
}

export default App;
