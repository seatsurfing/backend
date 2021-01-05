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
import { Ajax, Settings as OrgSettings, User } from 'flexspace-commons';
import Login from './pages/Login';
import LoginSuccess from './pages/LoginSuccess';
import LoginFailed from './pages/LoginFailed';
import ProtectedRoute from './pages/ProtectedRoute';
import Search from './pages/Search';
import SearchResult from './pages/SearchResult';
import Bookings from './pages/Bookings';
import ConfluenceHint from './pages/ConfluenceHint';
import RuntimeConfig from './components/RuntimeConfig';
import { AuthContext, AuthContextData } from './AuthContextData';
import Loading from './components/Loading';
import ConfluenceAnonymous from './pages/ConfluenceAnonymous';

interface Props {
}

class App extends React.Component<Props, AuthContextData> {
  constructor(props: Props) {
    super(props);
    this.state = {
      token: "",
      username: "",
      isLoading: true,
      maxBookingsPerUser: 0,
      maxDaysInAdvance: 0,
      maxBookingDurationHours: 0,
      setDetails: this.setDetails
    };
    if (window.location.href.indexOf("http://localhost") > -1 || window.location.href.indexOf("http://192.168.") > -1) {
      Ajax.DEV_MODE = true;
      Ajax.DEV_URL = "http://" + window.location.host.split(':').shift() + ":8090";
    }
    if (window.location.href.indexOf(".loca.lt/") > -1) {
      Ajax.DEV_MODE = true;
      Ajax.DEV_URL = "https://" + window.location.host.split(':').shift();
    }
    setTimeout(() => {
      this.verifyToken();
    }, 10);
  }

  verifyToken = async () => {
    let token: string | null = null;
    try {
      token = window.sessionStorage.getItem("jwt");
    } catch (e) {
      // Do nothing
    }
    if (token != null) {
      Ajax.JWT = token;
      User.getSelf().then(user => {
        this.loadSettings().then(() => {
          this.setDetails(token != null ? token : "", user.email);
          this.setState({ isLoading: false });
        });
      }).catch((e) => {
        Ajax.JWT = "";
        window.sessionStorage.removeItem("jwt");
        this.setState({ isLoading: false });
      });
    } else {
      this.setState({ isLoading: false });
    }
  }

  loadSettings = async () => {
    OrgSettings.list().then(settings => {
      let state: any = {};
      settings.forEach(s => {
        if (s.name === "max_bookings_per_user") state.maxBookingsPerUser = window.parseInt(s.value);
        if (s.name === "max_days_in_advance") state.maxDaysInAdvance = window.parseInt(s.value);
        if (s.name === "max_booking_duration_hours") state.maxBookingDurationHours = window.parseInt(s.value);
      });
      this.setState({
        ...this.state,
        ...state
      });
    });
  }

  setDetails = (token: string, username: string) => {
    this.loadSettings().then(() => {
      this.setState({
        token: token,
        username: username
      });
    });
  }

  render() {
    if (window !== window.parent) {
      // Add Confluence JS
      const script = document.createElement("script");
      script.src = "https://connect-cdn.atl-paas.net/all.js";
      document.head.appendChild(script);
      RuntimeConfig.EMBEDDED = true;
    }

    if (this.state.isLoading) {
      return <Loading />;
    }

    return (
      <Router basename={process.env.PUBLIC_URL}>
        <AuthContext.Provider value={this.state}>
          <Switch>
            <Route path="/login/confluence/anonymous" component={ConfluenceAnonymous} />
            <Route path="/login/confluence/:id" component={ConfluenceHint} />
            <Route path="/login/success/:id" component={LoginSuccess} />
            <Route path="/login/failed" component={LoginFailed} />
            <Route path="/login" component={Login} />
            <ProtectedRoute path="/search/result" component={SearchResult} />
            <ProtectedRoute path="/search" component={Search} />
            <ProtectedRoute path="/bookings" component={Bookings} />
            <Route path="/"><Redirect to="/login" /></Route>
          </Switch>
        </AuthContext.Provider>
      </Router>
    );
  }
}

export default withTranslation()(App as any);
