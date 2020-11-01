import React from 'react';
import './Login.css';
import { Form, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';

interface State {
}

export default class LoginFailed extends React.Component<{}, State> {
  render() {
    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Alert variant="danger">Login fehlgeschlagen.</Alert>
          <p>Möglicherweise ist das verwendete Konto nicht für diese Organisation freigeschaltet.</p>
          <Link className="btn btn-primary" to="/login">Zurück</Link>
        </Form>
      </div>
    )
  }
}
