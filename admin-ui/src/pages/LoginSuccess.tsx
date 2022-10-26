import React from 'react';
import './Login.css';
import Loading from '../components/Loading';
import { Form } from 'react-bootstrap';
import { Ajax, JwtDecoder, User } from 'flexspace-commons';
import { Navigate, Params, PathRouteProps } from 'react-router-dom';
import { withRouter } from '../types/withRouter';

interface State {
  redirect: string | null
}

interface Props extends PathRouteProps {
  params: Readonly<Params<string>>
  id: string
}

class LoginSuccess extends React.Component<Props, State> {
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
    if (this.props.params.id) {
      return Ajax.get("/auth/verify/" + this.props.params.id).then(res => {
        if (res.json && res.json.accessToken) {
          let jwtPayload = JwtDecoder.getPayload(res.json.accessToken);
          if (jwtPayload.role < User.UserRoleSpaceAdmin) {
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
          if (res.json.longLived) {
            Ajax.PERSISTER.persistRefreshTokenInLocalStorage(Ajax.CREDENTIALS);
          }
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
      return <Navigate replace={true} to={this.state.redirect} />
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

export default withRouter(LoginSuccess as any);
