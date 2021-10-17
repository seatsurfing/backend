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
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Locations from './pages/Locations';
import EditLocation from './pages/EditLocation';
import EditAuthProvider from './pages/EditAuthProvider';
import LoginSuccess from './pages/LoginSuccess';
import LoginFailed from './pages/LoginFailed';
import ProtectedRoute from './pages/ProtectedRoute';
import { Ajax } from 'flexspace-commons';
import Users from './pages/Users';
import EditUser from './pages/EditUser';
import Settings from './pages/Settings';
import Bookings from './pages/Bookings';
import SearchResult from './pages/SearchResult';
import ConfirmSignup from './pages/ConfirmSignup';
import Organizations from './pages/Organizations';
import EditOrganization from './pages/EditOrganization';
import Loading from './components/Loading';

interface Props {
}

interface State {
  isLoading: boolean
}

class App extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      isLoading: true
    };
    if (process.env.NODE_ENV.toLowerCase() === "development") {
      Ajax.URL = "http://" + window.location.host.split(':').shift() + ":8080";
    }
    setTimeout(() => {
      this.initAjax();
    }, 10);
  }

  initAjax = async () => {
    Ajax.PERSISTER.readCredentialsFromSessionStorage().then(c => {
      Ajax.CREDENTIALS = c;
      this.setState({
        isLoading: false
      });
    });
  }

  render() {
    if (this.state.isLoading) {
      return <Loading />;
    }

    return (
        <Router basename={process.env.PUBLIC_URL}>
          <Switch>
            <Route path="/login/success/:id" component={LoginSuccess} />
            <Route path="/login/failed" component={LoginFailed} />
            <Route path="/login" component={Login} />
            <Route path="/confirm/:id" component={ConfirmSignup} />
            <ProtectedRoute path="/dashboard" component={Dashboard} />
            <ProtectedRoute path="/locations/add" component={EditLocation} />
            <ProtectedRoute path="/locations/:id" component={EditLocation} />
            <ProtectedRoute path="/locations" component={Locations} />
            <ProtectedRoute path="/users/add" component={EditUser} />
            <ProtectedRoute path="/users/:id" component={EditUser} />
            <ProtectedRoute path="/users" component={Users} />
            <ProtectedRoute path="/settings/auth-providers/add" component={EditAuthProvider} />
            <ProtectedRoute path="/settings/auth-providers/:id" component={EditAuthProvider} />
            <ProtectedRoute path="/settings" component={Settings} />
            <ProtectedRoute path="/bookings" component={Bookings} />
            <ProtectedRoute path="/organizations/add" component={EditOrganization} />
            <ProtectedRoute path="/organizations/:id" component={EditOrganization} />
            <ProtectedRoute path="/organizations" component={Organizations} />
            <ProtectedRoute path="/search/:keyword" component={SearchResult} />
            <Route path="/"><Redirect to="/login" /></Route>
          </Switch>
        </Router>
    );
  }
}

export default withTranslation()(App as any);
