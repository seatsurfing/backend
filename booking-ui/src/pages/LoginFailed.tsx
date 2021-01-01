import React from 'react';
import './Login.css';
import { Form, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import RuntimeConfig from '../components/RuntimeConfig';

interface State {
}

interface Props {
  t: TFunction
}

class LoginFailed extends React.Component<Props, State> {
  render() {
    let backButton = <></>;
    if (!RuntimeConfig.EMBEDDED) {
      backButton = <Link className="btn btn-primary" to="/login">{this.props.t("back")}</Link>;
    }

    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Alert variant="danger">{this.props.t("errorLoginFailed")}</Alert>
          <p>{this.props.t("loginFailedDescription")}</p>
          {backButton}
        </Form>
      </div>
    )
  }
}

export default withTranslation()(LoginFailed as any);
