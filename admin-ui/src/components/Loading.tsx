import { WithTranslation, withTranslation } from 'next-i18next';
import React from 'react';
import { Loader as IconLoad } from 'react-feather';

interface State {
}

interface Props extends WithTranslation {
}

class Loading extends React.Component<Props, State> {
    render() {
        return (
            <div className="padding-top center"><IconLoad className="feather loader" /> {this.props.t("loadingHint")}</div>
        );
    }
}

export default withTranslation()(Loading as any);
