import React from 'react';
import './Login.css';
import Loading from '../components/Loading';
import { Form } from 'react-bootstrap';
import { Ajax } from 'flexspace-commons';
import RuntimeConfig from '../components/RuntimeConfig';
import { AuthContext } from '../AuthContextData';
import { Navigate, Params } from 'react-router-dom';
import { withRouter } from '../types/withRouter';

interface State {
  redirect: string | null
}

interface Props {
  params: Readonly<Params<string>>
}

class LoginSuccess extends React.Component<Props, State> {
  static contextType = AuthContext;
  
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
          Ajax.CREDENTIALS = {
            accessToken: res.json.accessToken,
            refreshToken: res.json.refreshToken,
            accessTokenExpiry: new Date(new Date().getTime() + Ajax.ACCESS_TOKEN_EXPIRY_OFFSET)
          };
          Ajax.PERSISTER.updateCredentialsSessionStorage(Ajax.CREDENTIALS).then(() => {
            RuntimeConfig.setLoginDetails(this.context).then(() => {
              this.setState({
                redirect: "/search"
              });
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
