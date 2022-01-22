import React from 'react';
import './Login.css';
import { Form, Alert } from 'react-bootstrap';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { RouteProps } from 'react-router-dom';

interface State {
}

interface Props extends RouteProps {
  t: TFunction
}

class ConfluenceAnonymous extends React.Component<Props, State> {
  render() {
    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Alert variant="danger">{this.props.t("errorConfluenceAnonymous")}</Alert>
          <p>{this.props.t("confluenceAnonymousHint")}</p>
        </Form>
      </div>
    )
  }
}

export default withTranslation()(ConfluenceAnonymous as any);
