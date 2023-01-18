import React from 'react';
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes
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
import { Ajax } from 'flexspace-commons';
import Users from './pages/Users';
import EditUser from './pages/EditUser';
import Settings from './pages/Settings';
import Bookings from './pages/Bookings';
import EditBooking from './pages/EditBooking';
import SearchResult from './pages/SearchResult';
import ConfirmSignup from './pages/ConfirmSignup';
import Organizations from './pages/Organizations';
import EditOrganization from './pages/EditOrganization';
import Loading from './components/Loading';
import ReportAnalysis from './pages/ReportAnalysis';
import { ProtectedRoute } from './pages/ProtectedRoute';

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
        <BrowserRouter basename={process.env.PUBLIC_URL}>
          <Routes>
            <Route path="/login/success/:id" element={<LoginSuccess />} />
            <Route path="/login/failed" element={<LoginFailed />} />
            <Route path="/login" element={<Login />} />
            <Route path="/confirm/:id" element={<ConfirmSignup />} />

            <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
            <Route path="/locations/add" element={<ProtectedRoute><EditLocation /></ProtectedRoute>} />
            <Route path="/locations/:id" element={<ProtectedRoute><EditLocation /></ProtectedRoute>} />
            <Route path="/locations" element={<ProtectedRoute><Locations /></ProtectedRoute>} />
            <Route path="/users/add" element={<ProtectedRoute><EditUser /></ProtectedRoute>} />
            <Route path="/users/:id" element={<ProtectedRoute><EditUser /></ProtectedRoute>} />
            <Route path="/users" element={<ProtectedRoute><Users /></ProtectedRoute>} />
            <Route path="/settings/auth-providers/add" element={<ProtectedRoute><EditAuthProvider /></ProtectedRoute>} />
            <Route path="/settings/auth-providers/:id" element={<ProtectedRoute><EditAuthProvider /></ProtectedRoute>} />
            <Route path="/settings" element={<ProtectedRoute><Settings /></ProtectedRoute>} />
            
            <Route path="/bookings/add" element={<ProtectedRoute><EditBooking /></ProtectedRoute>} />
            <Route path="/bookings/:id" element={<ProtectedRoute><EditBooking /></ProtectedRoute>} />
            <Route path="/bookings" element={<ProtectedRoute><Bookings /></ProtectedRoute>} />
            <Route path="/report/analysis" element={<ProtectedRoute><ReportAnalysis /></ProtectedRoute>} />
            <Route path="/organizations/add" element={<ProtectedRoute><EditOrganization /></ProtectedRoute>} />
            <Route path="/organizations/:id" element={<ProtectedRoute><EditOrganization /></ProtectedRoute>} />
            <Route path="/organizations" element={<ProtectedRoute><Organizations /></ProtectedRoute>} />
            <Route path="/search/:keyword" element={<ProtectedRoute><SearchResult /></ProtectedRoute>} />

            <Route path="/" element={<Navigate to="/login" />} />
            <Route path="*" element={<Navigate to="/login" />} />
          </Routes>
        </BrowserRouter>
    );
  }
}

export default withTranslation()(App as any);
