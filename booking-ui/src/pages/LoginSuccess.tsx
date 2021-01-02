import React from 'react';
import {
  RouteChildrenProps, Redirect
} from "react-router-dom";
import './Login.css';
import Loading from '../components/Loading';
import { Form } from 'react-bootstrap';
import { Ajax } from 'flexspace-commons';
import RuntimeConfig from '../components/RuntimeConfig';
import { AuthContext } from '../AuthContextData';

interface State {
  redirect: string | null
}

interface Props {
  id: string
}

export default class LoginSuccess extends React.Component<RouteChildrenProps<Props>, State> {
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
    if (this.props.match?.params.id) {
      return Ajax.get("/auth/verify/" + this.props.match.params.id).then(result => {
        if (result.json && result.json.jwt) {
          RuntimeConfig.setLoginDetails(result.json.jwt, this.context).then(() => {
            this.setState({
              redirect: "/search"
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
