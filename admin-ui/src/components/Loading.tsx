import { WithTranslation, withTranslation } from 'next-i18next';
import React from 'react';
import { Loader as IconLoad } from 'react-feather';

interface State {
}

interface Props extends WithTranslation {
}

class Loading extends React.Component<Props, State> {
    render() {
        let text = "Loading...";
        return (
            <div className="padding-top center"><IconLoad className="feather loader" /> {text}</div>
        );
    }
}

export default withTranslation(['admin'])(Loading as any);
