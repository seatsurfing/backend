import React from 'react';
import './CenterContent.css';
import { Loader as IconLoad } from 'react-feather';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { Link, Params, PathRouteProps } from 'react-router-dom';
import { Ajax } from 'flexspace-commons';
import { withRouter } from '../types/withRouter';

interface State {
  loading: boolean
  success: boolean
}

interface Props extends PathRouteProps {
  params: Readonly<Params<string>>
  t: TFunction
}

class ConfirmSignup extends React.Component<Props, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      loading: true,
      success: false
    };
  }

  componentDidMount = () => {
    this.loadData();
  }

  loadData = (id?: string) => {
    if (!id) {
      id = this.props.params.id;
    }
    if (id) {
      Ajax.postData("/signup/confirm/" + id, null).then((res) => {
        if (res.status >= 200 && res.status <= 299) {
          this.setState({ loading: false, success: true });
        } else {
          this.setState({ loading: false, success: false });
        }
      }).catch((e) => {
        this.setState({ loading: false, success: false });
      });
    } else {
      this.setState({ loading: false, success: false });
    }
  }

  render() {
    let loading = <></>;
    let result = <></>;
    if (this.state.loading) {
      loading = <div><IconLoad className="feather loader" /> {this.props.t("loadingHint")}</div>;
    } else {
      if (this.state.success) {
        result = (
          <div>
            <p>{this.props.t("orgSignupSuccess")}</p>
            <Link to="/login" className="btn btn-primary">{this.props.t("orgSignupGoToLogin")}</Link>
          </div>
        );
      } else {
        result = (
          <div>
            <p>{this.props.t("orgSignupFailed")}</p>
          </div>
        );
      }
    }

    return (
      <div className="container-center">
        <div className="container-center-inner">
          <img src="./seatsurfing.svg" alt="Seatsurfing" className="logo" />
          {loading}
          {result}
        </div>
      </div>
    );
  }
}

export default withRouter(withTranslation()(ConfirmSignup as any));
