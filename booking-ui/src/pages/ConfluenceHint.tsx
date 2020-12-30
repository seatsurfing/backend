import React from 'react';
import './Login.css';
import { Form, Alert, Button } from 'react-bootstrap';
import { RouteChildrenProps } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
}

interface RoutedProps {
  id: string
}

interface Props extends RouteChildrenProps<RoutedProps> {
  t: TFunction
}

class ConfluenceHint extends React.Component<Props, State> {
  onCreateAccountClick = () => {
    window.open("https://seatsurfing.de/");
  }

  render() {
    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Alert variant="danger">{this.props.t("errorConfluenceClientIdUnknown")}</Alert>
          <p>{this.props.t("confluenceClientIdHint")}</p>
          <pre>{this.props.match?.params.id}</pre>
          <p>{this.props.t("confluenceClientIdHint2")}</p>
          <Button className="btn btn-primary" onClick={this.onCreateAccountClick}>{this.props.t("createAccount")}</Button>
        </Form>
      </div>
    )
  }
}

export default withTranslation()(ConfluenceHint as any);
