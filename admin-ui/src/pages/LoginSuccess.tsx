import React from 'react';
import {
  RouteChildrenProps, Redirect
} from "react-router-dom";
import './Login.css';
import Loading from '../components/Loading';
import { Form } from 'react-bootstrap';
import { Ajax, JwtDecoder } from 'flexspace-commons';

interface State {
  redirect: string | null
}

interface Props {
  id: string
}

export default class LoginSuccess extends React.Component<RouteChildrenProps<Props>, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      redirect: null
    };
  }

  componentDidMount = () => {
    this.loadData();
  }

  loadData = () => {
    if (this.props.match?.params.id) {
      return Ajax.get("/auth/verify/" + this.props.match.params.id).then(res => {
        if (res.json && res.json.accessToken) {
          let jwtPayload = JwtDecoder.getPayload(res.json.accessToken);
          if (!jwtPayload.admin) {
            this.setState({
              redirect: "/login/failed"
            });
            return;
          }
          Ajax.CREDENTIALS = {
            accessToken: res.json.accessToken,
            refreshToken: res.json.refreshToken,
            accessTokenExpiry: new Date(new Date().getTime() + Ajax.ACCESS_TOKEN_EXPIRY_OFFSET)
          };
          Ajax.PERSISTER.updateCredentialsSessionStorage(Ajax.CREDENTIALS).then(() => {
            this.setState({
              redirect: "/dashboard"
            });
          });
        } else {
          this.setState({
            redirect: "/login/failed"
          });
        }
      }).catch(() => {
        this.setState({
          redirect: "/login/failed"
        });
      });
    }
  }

  render() {
    if (this.state.redirect != null) {
      return <Redirect to={this.state.redirect} />
    }

    return (
      <div className="container-signin">
        <Form className="form-signin">
          <Loading />
        </Form>
      </div>
    );
  }
}
