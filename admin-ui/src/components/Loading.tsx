import React from 'react';
import { Loader as IconLoad } from 'react-feather';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
}

interface Props {
    t: TFunction
}

class Loading extends React.Component<Props, State> {
    render() {
        return (
            <div className="padding-top center"><IconLoad className="feather loader" /> {this.props.t("loadingHint")}</div>
        );
    }
}

export default withTranslation()(Loading as any);
