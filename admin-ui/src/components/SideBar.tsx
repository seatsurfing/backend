import React from 'react';
import { Home as IconHome, Users as IconUsers, Map as IconMap, Book as IconBook, Settings as IconSettings, Box as IconBox, Activity as IconAnalysis, ExternalLink as IconExternalLink } from 'react-feather';
import { User } from 'flexspace-commons';
import { WithTranslation, withTranslation } from 'next-i18next';
import { Nav } from 'react-bootstrap';
import { NextRouter, withRouter } from 'next/router';

interface State {
    superAdmin: boolean
    spaceAdmin: boolean
    orgAdmin: boolean
}

interface Props extends WithTranslation {
    router: NextRouter
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

    getActiveKey = () => {
        const startPaths = [
            '/organizations',
            '/users',
            '/settings',
            '/locations',
            '/bookings'
        ];
        let path = this.props.router.pathname;
        let result = path;
        startPaths.forEach(startPath => {
            if (path.startsWith(startPath)) {
                result = startPath;
            }
        });
        return result;
    }

    render() {
        let orgItem = <></>;
        if (this.state.superAdmin) {
            orgItem = (
                <li className="nav-item">
                    <Nav.Link eventKey="/organizations" onClick={() => this.props.router.push("/organizations")}><IconBox className="feather" /> {this.props.t("organizations")}</Nav.Link>
                </li>
            );
        }
        let orgAdminItems = <></>;
        if (this.state.orgAdmin) {
            orgAdminItems = (
                <>
                    <li className="nav-item">
                        <Nav.Link eventKey="/users" onClick={() => this.props.router.push("/users")}><IconUsers className="feather" /> {this.props.t("users")}</Nav.Link>
                    </li>
                    <li className="nav-item">
                        <Nav.Link eventKey="/settings" onClick={() => this.props.router.push("/settings")}><IconSettings className="feather" /> {this.props.t("settings")}</Nav.Link>
                    </li>
                </>
            );
        }
        return (
            <Nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" activeKey={this.getActiveKey()}>
                <div className="sidebar-sticky pt-3">
                    <ul className="nav flex-column">
                        <li className="nav-item">
                            <Nav.Link eventKey="/dashboard" onClick={() => this.props.router.push("/dashboard")}><IconHome className="feather" /> {this.props.t("dashboard")}</Nav.Link>
                        </li>
                        <li className="nav-item">
                            <Nav.Link eventKey="/locations" onClick={() => this.props.router.push("/locations")}><IconMap className="feather" /> {this.props.t("areas")}</Nav.Link>
                        </li>
                        <li className="nav-item">
                            <Nav.Link eventKey="/bookings" onClick={() => this.props.router.push("/bookings")}><IconBook className="feather" /> {this.props.t("bookings")}</Nav.Link>
                        </li>
                        <li className="nav-item">
                            <Nav.Link eventKey="/report/analysis" onClick={() => this.props.router.push("/report/analysis")}><IconAnalysis className="feather" /> {this.props.t("analysis")}</Nav.Link>
                        </li>
                        {orgAdminItems}
                        {orgItem}
                        <li className="nav-item">
                            <Nav.Link onClick={(e) => {e.preventDefault(); window.location.href="/ui/";}}><IconExternalLink className="feather" /> {this.props.t("bookingui")}</Nav.Link>
                        </li>
                    </ul>
                </div>
            </Nav>
        );
    }
}

export default withTranslation()(withRouter(SideBar as any));
