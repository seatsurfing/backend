import React from 'react';
import './CenterContent.css';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { Ajax } from 'flexspace-commons';
import { Button, Form } from 'react-bootstrap';
import { Link } from 'react-router-dom';

interface State {
  loading: boolean
  complete: boolean
  success: boolean
  email: string
}

interface Props {
  t: TFunction
}

class InitPasswordReset extends React.Component<Props, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      loading: false,
      complete: false,
      success: false,
      email: ""
    };
  }

  onPasswordSubmit = (e: any) => {
    e.preventDefault();
    this.setState({ loading: true, complete: false, success: false });
    let payload = {
      "email": this.state.email
    };
    Ajax.postData("/auth/initpwreset", payload).then((res) => {
      if (res.status >= 200 && res.status <= 299) {
        this.setState({ loading: false, complete: true, success: true });
      } else {
        this.setState({ loading: false, complete: true, success: false });
      }
    }).catch((e) => {
      this.setState({ loading: false, complete: true, success: false });
    });
  }

  render() {
    if (this.state.complete) {
      if (this.state.success) {
        return (
          <div className="container-center">
            <div className="container-center-inner">
              <img src="./seatsurfing.svg" alt="Seatsurfing" className="logo" />
              <p>{this.props.t("initPasswordResetEmail")}</p>
            </div>
          </div>
        );
      } else {
        return (
          <div className="container-center">
            <div className="container-center-inner">
              <img src="./seatsurfing.svg" alt="Seatsurfing" className="logo" />
              <p>{this.props.t("initPasswordResetFailed")}</p>
            </div>
          </div>
        );
      }
    }

    return (
      <div className="container-center">
        <Form className="container-center-inner" onSubmit={this.onPasswordSubmit}>
          <img src="./seatsurfing.svg" alt="Seatsurfing" className="logo" />
          <Form.Group>
            <Form.Control type="email" placeholder={this.props.t("emailPlaceholder")} value={this.state.email} onChange={(e: any) => this.setState({ email: e.target.value })} required={true} autoFocus={true} />
          </Form.Group>
          <Button className="margin-top-10" variant="primary" type="submit" disabled={this.state.loading}>{this.props.t("changePassword")}</Button>
          <p className="margin-top-50"><Link to="/login">{this.props.t("back")}</Link></p>
        </Form>
      </div>
    );
  }
}

export default withTranslation()(InitPasswordReset as any);
