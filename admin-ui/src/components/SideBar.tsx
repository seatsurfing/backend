import React from 'react';
import './SideBar.css';
import { Home as IconHome, Users as IconUsers, Map as IconMap, Book as IconBook, Settings as IconSettings, Box as IconBox } from 'react-feather';
import { NavLink } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { User } from 'flexspace-commons';

interface State {
    superAdmin: boolean
}

interface Props {
    t: TFunction
}

class SideBar extends React.Component<Props, State> {
    constructor(props: any) {
        super(props);
        this.state = {
            superAdmin: false
        };
    }

    componentDidMount = () => {
        User.getSelf().then(user => {
            this.setState({superAdmin: user.superAdmin});
        });
    }

    render() {
        let orgItem = <></>;
        if (this.state.superAdmin) {
            orgItem = (
                <li className="nav-item">
                    <NavLink to="/organizations" className="nav-link" activeClassName="active"><IconBox className="feather" /> {this.props.t("organizations")}</NavLink>
                </li>
            );
        }
        return (
            <nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div className="sidebar-sticky pt-3">
                    <ul className="nav flex-column">
                        <li className="nav-item">
                            <NavLink to="/dashboard" className="nav-link" activeClassName="active"><IconHome className="feather" /> {this.props.t("dashboard")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/locations" className="nav-link" activeClassName="active"><IconMap className="feather" /> {this.props.t("areas")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/users" className="nav-link" activeClassName="active"><IconUsers className="feather" /> {this.props.t("users")}</NavLink>
                        </li>
                        {orgItem}
                        <li className="nav-item">
                            <NavLink to="/bookings" className="nav-link" activeClassName="active"><IconBook className="feather" /> {this.props.t("bookings")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/settings" className="nav-link" activeClassName="active"><IconSettings className="feather" /> {this.props.t("settings")}</NavLink>
                        </li>
                    </ul>
                </div>
            </nav>
        );
    }
}

export default withTranslation()(SideBar as any);
