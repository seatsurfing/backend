import React from 'react';
import './SideBar.css';
import { Home as IconHome, Users as IconUsers, Map as IconMap, Book as IconBook, Settings as IconSettings, Box as IconBox, Activity as IconAnalysis, ExternalLink as IconExternalLink } from 'react-feather';
import { NavLink } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { User } from 'flexspace-commons';

interface State {
    superAdmin: boolean
    spaceAdmin: boolean
    orgAdmin: boolean
}

interface Props {
    t: TFunction
}

class SideBar extends React.Component<Props, State> {
    constructor(props: any) {
        super(props);
        this.state = {
            superAdmin: false,
            spaceAdmin: false,
            orgAdmin: false,
        };
    }

    componentDidMount = () => {
        User.getSelf().then(user => {
            this.setState({
                superAdmin: user.superAdmin,
                spaceAdmin: user.spaceAdmin,
                orgAdmin: user.admin,
            });
        });
    }

    render() {
        let orgItem = <></>;
        if (this.state.superAdmin) {
            orgItem = (
                <li className="nav-item">
                    <NavLink to="/organizations" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconBox className="feather" /> {this.props.t("organizations")}</NavLink>
                </li>
            );
        }
        let orgAdminItems = <></>;
        if (this.state.orgAdmin) {
            orgAdminItems = (
                <>
                    <li className="nav-item">
                        <NavLink to="/users" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconUsers className="feather" /> {this.props.t("users")}</NavLink>
                    </li>
                    <li className="nav-item">
                        <NavLink to="/settings" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconSettings className="feather" /> {this.props.t("settings")}</NavLink>
                    </li>
                </>
            );
        }
        return (
            <nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div className="sidebar-sticky pt-3">
                    <ul className="nav flex-column">
                        <li className="nav-item">
                            <NavLink to="/dashboard" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconHome className="feather" /> {this.props.t("dashboard")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/locations" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconMap className="feather" /> {this.props.t("areas")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/bookings" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconBook className="feather" /> {this.props.t("bookings")}</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/report/analysis" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconAnalysis className="feather" /> {this.props.t("analysis")}</NavLink>
                        </li>
                        {orgAdminItems}
                        {orgItem}
                        <li className="nav-item">
                            <NavLink to="/bookingui" onClick={() => {window.location.href="/ui/"}} className={({isActive}) => "nav-link " + (isActive ? "active" : "")}><IconExternalLink className="feather" /> {this.props.t("bookingui")}</NavLink>
                        </li>
                    </ul>
                </div>
            </nav>
        );
    }
}

export default withTranslation()(SideBar as any);
