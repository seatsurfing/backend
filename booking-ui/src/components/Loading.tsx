import { WithTranslation, withTranslation } from 'next-i18next';
import React from 'react';
import { Loader as IconLoad } from 'react-feather';

interface State {
}

interface Props extends WithTranslation {
    showText: boolean
    paddingTop: boolean
}

class Loading extends React.Component<Props, State> {
    render() {
        return (
            <div className={this.props.paddingTop === undefined || this.props.paddingTop === true ? "padding-top center" : "center"}><IconLoad className="feather loader" />{this.props.showText === undefined || this.props.showText === true ? " " + this.props.t("loadingHint") : ""}</div>
        );
    }
}

export default withTranslation()(Loading as any);
