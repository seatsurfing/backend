import React from 'react';
import './SideBar.css';
import { Home as IconHome, Users as IconUsers, Map as IconMap, Book as IconBook, Settings as SettingsIcon } from 'react-feather';
import { NavLink } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
}

interface Props {
    t: TFunction
}

class SideBar extends React.Component<Props, State> {
    render() {
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
                        <li className="nav-item">
                            <NavLink to="/bookings" className="nav-link" activeClassName="active"><IconBook className="feather" /> {this.props.t("bookings")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/settings" className="nav-link" activeClassName="active"><SettingsIcon className="feather" /> {this.props.t("settings")}</NavLink>
                        </li>
                    </ul>
                </div>
            </nav>
        );
    }
}

export default withTranslation()(SideBar as any);
