import React from 'react';
import './Login.css';
import { Route, Redirect } from 'react-router-dom';
import NavBar from '../components/NavBar';

export default class ProtectedRoute extends Route {
  render() {
    let token: string | null = null;
    try {
      token = window.sessionStorage.getItem("jwt");
    } catch (e) {
      // Do nothing
    }
    if (!token) {
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
