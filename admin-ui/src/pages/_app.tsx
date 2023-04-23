import 'bootstrap/dist/css/bootstrap.min.css'
import '@/styles/App.css'
import '@/styles/CenterContent.css'
import '@/styles/Dashboard.css'
import '@/styles/EditLocation.css'
import '@/styles/Login.css'
import '@/styles/NavBar.css'
import '@/styles/Settings.css'
import '@/styles/SideBar.css'
import type { AppProps } from 'next/app'
import nextI18nConfig from '../../next-i18next.config'
import { WithTranslation, appWithTranslation } from 'next-i18next'
import { Ajax } from 'flexspace-commons'
import React from 'react'
import Loading from '@/components/Loading'
import dynamic from 'next/dynamic'
import Head from 'next/head'

interface State {
  isLoading: boolean;
}

interface Props extends WithTranslation, AppProps {
}

class App extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      isLoading: true
    };
    if (typeof window !== 'undefined') {
      if (process.env.NODE_ENV.toLowerCase() === "development") {
        Ajax.URL = "http://" + window.location.host.split(':').shift() + ":8080";
      }
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

    const { Component, pageProps } = this.props;
    return (
      <>
        <Head>
          <link rel="icon" href="/admin/favicon.ico" />
          <meta name="viewport" content="width=device-width, initial-scale=1" />
          <meta name="theme-color" content="#343a40" />
          <link rel="manifest" href="/admin/manifest.json" />
          <meta name="apple-mobile-web-app-capable" content="yes" />
          <meta name="apple-mobile-web-app-status-bar-style" content="default" />
          <link rel="shortcut icon" href="/admin/favicon-192.png" />
          <link rel="apple-touch-icon" href="/admin/favicon-192.png" />
          <link rel="apple-touch-startup-image" href="/admin/favicon-1024.png" />
          <title>Seatsurfing</title>
        </Head>
        <Component {...pageProps} />
      </>
    );
  }
}

/*
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
<Route path="/bookings" element={<ProtectedRoute><Bookings /></ProtectedRoute>} />
<Route path="/report/analysis" element={<ProtectedRoute><ReportAnalysis /></ProtectedRoute>} />
<Route path="/organizations/add" element={<ProtectedRoute><EditOrganization /></ProtectedRoute>} />
<Route path="/organizations/:id" element={<ProtectedRoute><EditOrganization /></ProtectedRoute>} />
<Route path="/organizations" element={<ProtectedRoute><Organizations /></ProtectedRoute>} />
<Route path="/search/:keyword" element={<ProtectedRoute><SearchResult /></ProtectedRoute>} />
*/

const NoSSRApp = dynamic(async () => App, { ssr: false });
export default appWithTranslation(NoSSRApp, nextI18nConfig);
