import React from 'react';
import './Login.css';
import { Form, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
}

interface Props {
  t: TFunction
}

class LoginFailed extends React.Component<Props, State> {
  render() {
    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Alert variant="danger">{this.props.t("errorLoginFailed")}</Alert>
          <p>{this.props.t("loginFailedDescription")}</p>
          <Link className="btn btn-primary" to="/login">{this.props.t("back")}</Link>
        </Form>
      </div>
    )
  }
}

export default withTranslation()(LoginFailed as any);
