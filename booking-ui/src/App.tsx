import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect
} from "react-router-dom";
import './i18n';
import './App.css';
import { withTranslation } from 'react-i18next';
import { Ajax } from 'flexspace-commons';
import Login from './pages/Login';
import LoginSuccess from './pages/LoginSuccess';
import LoginFailed from './pages/LoginFailed';
import ProtectedRoute from './pages/ProtectedRoute';
import Search from './pages/Search';
import SearchResult from './pages/SearchResult';
import Bookings from './pages/Bookings';
import ConfluenceHint from './pages/ConfluenceHint';
import RuntimeConfig from './components/RuntimeConfig';

interface Props {
}

class App extends React.Component<Props, {}> {
  render() {
    if (window !== window.parent) {
      // Add Confluence JS
      const script = document.createElement("script");
      script.src = "https://connect-cdn.atl-paas.net/all.js";
      document.head.appendChild(script);
      RuntimeConfig.EMBEDDED = true;
    }
    let jwt = window.sessionStorage.getItem("jwt");
    if (jwt) {
      Ajax.JWT = jwt;
    }
    if (window.location.href.indexOf("http://localhost") > -1 || window.location.href.indexOf("http://192.168.") > -1) {
      Ajax.DEV_MODE = true;
      Ajax.DEV_URL = "http://" + window.location.host.split(':').shift() + ":8090";
    }
    if (window.location.href.indexOf(".loca.lt/") > -1) {
      Ajax.DEV_MODE = true;
      Ajax.DEV_URL = "https://" + window.location.host.split(':').shift();
    }
    return (
      <Router basename={process.env.PUBLIC_URL}>
        <Switch>
          <Route path="/login/confluence/:id" component={ConfluenceHint} />
          <Route path="/login/success/:id" component={LoginSuccess} />
          <Route path="/login/failed" component={LoginFailed} />
          <Route path="/login" component={Login} />
          <ProtectedRoute path="/search/result" component={SearchResult} />
          <ProtectedRoute path="/search" component={Search} />
          <ProtectedRoute path="/bookings" component={Bookings} />
          <Route path="/"><Redirect to="/login" /></Route>
        </Switch>
      </Router>
    );
  }
}

export default withTranslation()(App as any);
