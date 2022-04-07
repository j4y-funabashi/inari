import {AuthState, onAuthUIStateChange} from '@aws-amplify/ui-components';
import {AmplifyAuthenticator, AmplifySignOut} from '@aws-amplify/ui-react';
import React from 'react';
import {BrowserRouter, Link, Route, Switch} from 'react-router-dom';
import MediaTimelinePage from './components/MediaTimelinePage';
import MediaTimelineMonthPage from './components/MediaTimelineMonthPage';
import MediaDetailPage from './components/MediaDetailPage';

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

  <BrowserRouter>


    <Switch>
      <Route exact path="/">
        <MediaTimelinePage />
      </Route>
      <Route path="/time/month/:monthid">
        <MediaTimelineMonthPage />
      </Route>
      <Route path="/media/:mediaid">
        <MediaDetailPage />
      </Route>
    </Switch>

    <div>
      <nav>
        <Link to="/">Timeline</Link>
        <div>Hello, {user.username}</div>
        <AmplifySignOut />
      </nav>
    </div>

  </BrowserRouter>

  ) : (
    <AmplifyAuthenticator />
  );
}

export default App;
