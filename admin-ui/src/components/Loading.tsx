import React from 'react';
import { Loader as IconLoad } from 'react-feather';
import './Loading.css';

export default class Loading extends React.Component {
    render() {
        return (
            <div><IconLoad className="feather loader" /> Daten werden geladen...</div>
        );
    }
}
