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
        let text = "";
        if ((this.props.showText === undefined) || (this.props.showText === true)) {
            if (this.props.tReady) {
                text = this.props.t("loadingHint");
            } else {
                text = " Loading...";
            }
        }
        return (
            <div className={this.props.paddingTop === undefined || this.props.paddingTop === true ? "padding-top center" : "center"}><IconLoad className="feather loader" />{text}</div>
        );
    }
}

export default withTranslation()(Loading as any);
