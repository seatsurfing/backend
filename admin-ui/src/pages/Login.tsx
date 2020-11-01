import React from 'react';
import { Form, Button } from 'react-bootstrap';
import {
  Redirect
} from "react-router-dom";
import './Login.css';
import { Organization, AuthProvider, Ajax, JwtDecoder } from 'flexspace-commons';

interface State {
  email: string
  password: string
  invalid: boolean
  redirect: string | null
  requirePassword: boolean
  providers: AuthProvider[] | null
}

export default class Login extends React.Component<{}, State> {
  org: Organization | null;

  constructor(props: any) {
    super(props);
    this.org = null;
    this.state = {
      email: "",
      password: "",
      invalid: false,
      redirect: null,
      requirePassword: false,
      providers: null
    };
  }

  onSubmit = (e: any) => {
    e.preventDefault();
    let email = this.state.email.split("@");
    if (email.length !== 2) {
      // Error
      return;
    }
    let payload = {
      email: this.state.email
    };
    Ajax.postData("/auth/preflight", payload).then((res) => {
      this.org = new Organization();
      this.org.deserialize(res.json.organization);
      this.setState({
        providers: res.json.authProviders,
        requirePassword: res.json.requirePassword
      });
    }).catch((e) => {
      this.setState({
        invalid: true
      });
    });
  }

  onPasswordSubmit = (e: any) => {
    e.preventDefault();
    let payload = {
      email: this.state.email,
      password: this.state.password
    };
    Ajax.postData("/auth/login", payload).then((res) => {
      let jwtPayload = JwtDecoder.getPayload(res.json.jwt);
      if (!jwtPayload.admin) {
        this.setState({
          invalid: true
        });
        return;
      }
      Ajax.JWT = res.json.jwt;
      window.sessionStorage.setItem("jwt", res.json.jwt);
      this.setState({
        redirect: "/dashboard"
      });
    }).catch((e) => {
      this.setState({
        invalid: true
      });
    });
  }

  cancelPasswordLogin = (e: any) => {
    e.preventDefault();
    this.setState({
      requirePassword: false,
      providers: null,
      invalid: false
    });
  }

  renderAuthProviderButton = (provider: AuthProvider) => {
    return (
      <p key={provider.id}>
        <Button variant="primary" className="btn-auth-provider" onClick={() => this.useProvider(provider)}>{provider.name}</Button>
      </p>
    );
  }

  useProvider = (provider: AuthProvider) => {
    let target = Ajax.getBackendUrl() + "/auth/" + provider.id + "/login/web";
    window.location.href = target;
  }

  render() {
    if (this.state.redirect != null) {
      return <Redirect to={this.state.redirect} />
    }

    if (this.state.requirePassword) {
      return (
        <div className="container-signin">
          <Form className="form-signin" onSubmit={this.onPasswordSubmit}>
            <p>Als {this.state.email} an {this.org?.name} anmelden:</p>
            <Form.Control type="password" placeholder="Kennwort" value={this.state.password} onChange={(e: any) => this.setState({ password: e.target.value, invalid: false })} required={true} isInvalid={this.state.invalid} minLength={8} autoFocus={true} />
            <Form.Control.Feedback type="invalid">Ungültiges Kennwort.</Form.Control.Feedback>
            <p><Button variant="primary" type="submit" className="btn-auth-provider">Anmelden</Button></p>
            <Button variant="secondary" className="btn-auth-provider" onClick={this.cancelPasswordLogin}>Zurück</Button>
          </Form>
        </div>
      );
    }

    if (this.state.providers != null) {
      let buttons = this.state.providers.map(provider => this.renderAuthProviderButton(provider));
      let providerSelection = <p>Als {this.state.email} an {this.org?.name} anmelden mit:</p>;
      if (buttons.length === 0) {
        providerSelection = <p>Für diesen Nutzer stehen keine Anmelde-Möglichkeiten zur Verfügung.</p>
      }
      return (
        <div className="container-signin">
          <Form className="form-signin">
            {providerSelection}
            {buttons}
            <Button variant="secondary" className="btn-auth-provider" onClick={() => this.setState({ providers: null })}>Zurück</Button>
          </Form>
        </div>
      );
    }

    return (
      <div className="container-signin">
        <Form className="form-signin" onSubmit={this.onSubmit}>
          <img src="./seatsurfing.svg" alt="Seatsurfing" className="logo" />
          <h3>Organisation verwalten.</h3>
          <Form.Control type="email" placeholder="E-Mail Adresse" value={this.state.email} onChange={(e: any) => this.setState({ email: e.target.value, invalid: false })} required={true} isInvalid={this.state.invalid} autoFocus={true} />
          <Form.Control.Feedback type="invalid">Ungültige E-Mail-Adresse.</Form.Control.Feedback>
          <Button variant="primary" type="submit">Anmelden</Button>
        </Form>
      </div>
    );
  }
}
