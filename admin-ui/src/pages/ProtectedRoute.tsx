import React from 'react';
import './Login.css';
import { Route, Redirect } from 'react-router-dom';
import { Ajax } from 'flexspace-commons';

export default class ProtectedRoute extends Route {
  render() {
    if (!Ajax.CREDENTIALS.accessToken) {
        return (
            <Redirect to="/login" />
        );
    }
    return (
        <Route path={this.props.path} component={this.props.component} />
    );
  }
}
