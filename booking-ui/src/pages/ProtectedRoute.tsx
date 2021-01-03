import React from 'react';
import './Login.css';
import { Route, Redirect } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { Ajax } from 'flexspace-commons';

export default class ProtectedRoute extends Route {
  render() {
    if (!Ajax.JWT) {
        return (
            <Redirect to="/login" />
        );
    }
    return (
        <>
          <NavBar />
          <Route path={this.props.path} component={this.props.component} />
        </>
    );
  }
}
