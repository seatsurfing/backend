import 'bootstrap/dist/css/bootstrap.min.css'
import '@/styles/App.css'
import '@/styles/NavBar.css'
import '@/styles/CenterContent.css'
import '@/styles/ConfluenceHint.css'
import '@/styles/Login.css'
import '@/styles/Search.css'
import type { AppProps } from 'next/app'
import nextI18nConfig from '../../next-i18next.config'
import { WithTranslation, appWithTranslation } from 'next-i18next'
import { Ajax, AjaxCredentials, User, Settings as OrgSettings } from 'flexspace-commons'
import RuntimeConfig from '@/components/RuntimeConfig'
import React from 'react'
//import { AuthContextData, AuthContextProvider } from '@/AuthContextData'
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
      isLoading: true,
      /*
      username: "",
      maxBookingsPerUser: 0,
      maxDaysInAdvance: 0,
      maxBookingDurationHours: 0,
      dailyBasisBooking: false,
      showNames: false,
      defaultTimezone: "",
      */
      //setDetails: this.setDetails
    };
    if (typeof window !== 'undefined') {
      if (process.env.NODE_ENV.toLowerCase() === "development") {
        Ajax.URL = "http://" + window.location.host.split(':').shift() + ":8080";
      }
      if (window.location.href.indexOf(".loca.lt/") > -1) {
        Ajax.URL = "https://" + window.location.host.split(':').shift();
      }
    }
    setTimeout(() => {
      RuntimeConfig.verifyToken(() => {
        this.setState({isLoading: false});
      });
    }, 10);
  }

  /*
  verifyToken = async () => {
    Ajax.CREDENTIALS = await Ajax.PERSISTER.readCredentialsFromSessionStorage();
    if (!Ajax.CREDENTIALS.accessToken) {
      Ajax.CREDENTIALS = await Ajax.PERSISTER.readRefreshTokenFromLocalStorage();
      if (Ajax.CREDENTIALS.refreshToken) {
        await Ajax.refreshAccessToken(Ajax.CREDENTIALS.refreshToken);
      }
    }
    if (Ajax.CREDENTIALS.accessToken) {
      User.getSelf().then(user => {
        this.loadSettings().then(() => {
          this.setDetails(user.email);
          this.setState({ isLoading: false });
        });
      }).catch((e) => {
        Ajax.CREDENTIALS = new AjaxCredentials();
        Ajax.PERSISTER.deleteCredentialsFromSessionStorage().then(() => {
          this.setState({ isLoading: false });
        });
      });
    } else {
      this.setState({ isLoading: false });
    }
  }

  loadSettings = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      OrgSettings.list().then(settings => {
        let state: any = {};
        settings.forEach(s => {
          if (typeof window !== 'undefined') {
            if (s.name === "max_bookings_per_user") state.maxBookingsPerUser = window.parseInt(s.value);
            if (s.name === "max_days_in_advance") state.maxDaysInAdvance = window.parseInt(s.value);
            if (s.name === "max_booking_duration_hours") state.maxBookingDurationHours = window.parseInt(s.value);
          }
          if (s.name === "daily_basis_booking") state.dailyBasisBooking = (s.value === "1");
          if (s.name === "show_names") state.showNames = (s.value === "1");
          if (s.name === "default_timezone") state.defaultTimezone = s.value;
        });
        self.setState({
          ...self.state,
          ...state
        }, () => resolve());
      });
    });
  }

  setDetails = (username: string) => {
    this.loadSettings().then(() => {
      this.setState({
        username: username
      });
    });
  }
  */

  render() {
    if (typeof window !== 'undefined') {
      if (window !== window.parent) {
        // Add Confluence JS
        if (!document.getElementById("confluence-js")) {
          const script = document.createElement("script");
          script.id = "confluence-js";
          script.src = "https://connect-cdn.atl-paas.net/all.js";
          document.head.appendChild(script);
        }
        RuntimeConfig.EMBEDDED = true;
      }
    }

    if (this.state.isLoading) {
      return <Loading />;
    }

    const { Component, pageProps } = this.props;
    return (
      <>
        <Head>
          <link rel="icon" href="/ui/favicon.ico" />
          <meta name="viewport" content="width=device-width, initial-scale=1" />
          <meta name="theme-color" content="#343a40" />
          <link rel="manifest" href="/ui/manifest.json" />
          <meta name="apple-mobile-web-app-capable" content="yes" />
          <meta name="apple-mobile-web-app-status-bar-style" content="default" />
          <link rel="shortcut icon" href="/ui/favicon-192.png" />
          <link rel="apple-touch-icon" href="/ui/favicon-192.png" />
          <link rel="apple-touch-startup-image" href="/ui/favicon-1024.png" />
          <title>Seatsurfing</title>
        </Head>
        <Component {...pageProps} />
      </>
    );
  }
  /*
  <AuthContextProvider>
    <Component {...pageProps} />
  </AuthContextProvider>
  */
}

const NoSSRApp = dynamic(async () => App, { ssr: false });
export default appWithTranslation(NoSSRApp, nextI18nConfig);


