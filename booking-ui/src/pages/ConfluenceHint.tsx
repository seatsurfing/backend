import React from 'react';
import './ConfluenceHint.css';
import { Button } from 'react-bootstrap';
import { Copy as IconCopy } from 'react-feather';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { Params, PathRouteProps } from 'react-router-dom';
import { withRouter } from '../types/withRouter';

interface State {
}

interface Props extends PathRouteProps {
  params: Readonly<Params<string>>
  t: TFunction
}

class ConfluenceHint extends React.Component<Props, State> {
  onCreateAccountClick = (e: any) => {
    e.preventDefault();
    window.open("https://seatsurfing.app/sign-up/");
  }

  onAdminClick = (e: any) => {
    e.preventDefault();
    window.open("https://app.seatsurfing.app/admin/");
  }

  onHelpClick = (e: any) => {
    e.preventDefault();
    window.open("https://seatsurfing.app/contact/");
  }

  onInputClick = (e: any) => {
    e.preventDefault();
    let input = document.querySelector("input.copy-input");
    if (input) {
      (input as HTMLInputElement).select();
    }
  }

  onCopyClick = () => {
    let input = document.querySelector("input.copy-input");
    if (input) {
      (input as HTMLInputElement).select();
      document.execCommand("copy");
    }
  }

  render() {
    return (
      <div className="container-confluence">
          <h1>{this.props.t("errorConfluenceClientIdUnknown")}</h1>
          <p>{this.props.t("confluenceClientIdHint")}</p>
          <ol>
            <li><Button variant="link" className="button-link" onClick={this.onCreateAccountClick}>{this.props.t("confluenceClientIdStep1")}</Button></li>
            <li><Button variant="link" className="button-link" onClick={this.onAdminClick}>{this.props.t("confluenceClientIdStep2")}</Button></li>
            <li>{this.props.t("confluenceClientIdStep3")}</li>
            <li>
              {this.props.t("confluenceClientIdStep4")}
              <br />
              <input type="text" className="copy-input" size={36} onClick={this.onInputClick} value={this.props.params.id} readOnly={true} />
              <Button variant="link" size="sm" className="copy-button" onClick={this.onCopyClick}><IconCopy className="feather" /></Button>
            </li>
            <li>{this.props.t("confluenceClientIdStep5")}</li>
            <li>{this.props.t("confluenceClientIdStep6")}</li>
          </ol>
          <p>{this.props.t("confluenceClientIdHint2")} <Button variant="link" className="button-link" onClick={this.onHelpClick}>{this.props.t("confluenceClientIdHint3")}</Button></p>
      </div>
    )
  }
}

export default withRouter(withTranslation()(ConfluenceHint as any));
