import React from 'react';
import { Navbar, Nav } from 'react-bootstrap';
import { NavLink, Redirect } from 'react-router-dom';
import { Ajax, JwtDecoder } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import './NavBar.css';

interface State {
    redirect: string | null
}

interface Props {
    t: TFunction
}

class NavBar extends React.Component<Props, State> {
    constructor(props: any) {
        super(props);
        this.state = {
            redirect: null
        };
    }

    logOut = (e: any) => {
        e.preventDefault();
        Ajax.JWT = "";
        window.sessionStorage.removeItem("jwt");
        this.setState({
            redirect: "/login"
        });
    }

    render() {
        if (this.state.redirect != null) {
            let target = this.state.redirect;
            this.setState({ redirect: null });
            return <Redirect to={target} />
        }

        let jwt = JwtDecoder.getPayload(Ajax.JWT);
        let username = jwt.email;

        return (
            <Navbar bg="dark" variant="dark" expand="lg" fixed="top">
                <Navbar.Brand as={NavLink} to="/search"><img src="/ui/seatsurfing_white.svg" alt="Seatsurfing" /></Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="mr-auto">
                        <Nav.Link as={NavLink} to="/search" activeClassName="active">{this.props.t("bookSeat")}</Nav.Link>
                        <Nav.Link as={NavLink} to="/bookings" activeClassName="active">{this.props.t("myBookings")}</Nav.Link>
                        <Nav.Link onClick={this.logOut}>{this.props.t("signout")}</Nav.Link>
                    </Nav>
                    <Nav className="mr-right">
                        <Navbar.Text>{username}</Navbar.Text>
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        );
    }
}

export default withTranslation()(NavBar as any);
