import React from 'react';
import { Navbar, Nav, Modal, Button, Form, Badge } from 'react-bootstrap';
import { Navigate, NavLink } from 'react-router-dom';
import { Ajax, User, MergeRequest, AjaxCredentials } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import './NavBar.css';
import RuntimeConfig from './RuntimeConfig';
import { Users as IconMerge, Bell as IconAlert } from 'react-feather';
import { AuthContext } from '../AuthContextData';

interface State {
    redirect: string | null
    showMergeInit: boolean
    showMergeNextStep: boolean
    showMergeRequests: boolean
    targetUserEmail: string
    invalidTargetUserEmail: boolean
    mergeRequests: MergeRequest[]
    allowMergeInit: boolean
}

interface Props {
    t: TFunction
}

class NavBar extends React.Component<Props, State> {
    static contextType = AuthContext;

    constructor(props: any) {
        super(props);
        this.state = {
            redirect: null,
            showMergeInit: false,
            showMergeNextStep: false,
            showMergeRequests: false,
            targetUserEmail: "",
            invalidTargetUserEmail: false,
            mergeRequests: [],
            allowMergeInit: false
        };
    }

    componentDidMount = () => {
        this.loadData();
    }

    loadData = () => {
        User.getMergeRequests().then(list => {
            this.setState({ mergeRequests: list });
        });
        User.getSelf().then(user => {
            if (user.email === user.atlassianId) {
                this.setState({ allowMergeInit: true });
            }
        });
    }

    logOut = (e: any) => {
        e.preventDefault();
        Ajax.CREDENTIALS = new AjaxCredentials();
        Ajax.PERSISTER.deleteCredentialsFromSessionStorage().then(() => {
            this.setState({
                redirect: "/login"
            });
        });
    }

    showMergeModal = (e: any) => {
        e.preventDefault();
        this.setState({ showMergeInit: true });
    }

    showMergeRequestsModal = (e: any) => {
        e.preventDefault();
        this.setState({ showMergeRequests: true });
    }

    initMerge = () => {
        User.initMerge(this.state.targetUserEmail).then(() => {
            this.setState({
                showMergeInit: false,
                showMergeNextStep: true
            })
        }).catch(() => {
            this.setState({
                invalidTargetUserEmail: true
            });
        });
    }

    openWebUI = () => {
        window.open(process.env.PUBLIC_URL);
    }

    acceptMergeRequest = (id: string) => {
        User.finishMerge(id).then(() => {
            this.setState({
                showMergeRequests: false
            });
            this.loadData();
        });
    }

    renderMergeRequest = (item: MergeRequest) => {
        return (
            <p key={item.id}>
                {item.email} <Button size="sm" onClick={() => this.acceptMergeRequest(item.id)}>{this.props.t("accept")}</Button>
            </p>
        );
    }

    render() {
        if (this.state.redirect != null) {
            let target = this.state.redirect;
            this.setState({ redirect: null });
            return <Navigate replace={true} to={target} />
        }

        let signOffButton = <></>;
        let initMergeButton = <></>;
        let mergeRequestsButton = <></>;
        let userInfo = <></>;
        let collapsable = <></>;
        let navClass = "navbar-expand";

        if (!RuntimeConfig.EMBEDDED) {
            navClass = "";
            signOffButton = <Nav.Link onClick={this.logOut}>{this.props.t("signout")}</Nav.Link>;
            userInfo = <Navbar.Text>{this.context.username}</Navbar.Text>;
            if (this.state.mergeRequests.length > 0) {
                mergeRequestsButton = <Nav.Link onClick={this.showMergeRequestsModal} className="icon-link"><IconAlert className="feather feather-lg" /><Badge pill={true} variant="light" className="badge-top">{this.state.mergeRequests.length}</Badge></Nav.Link>;
            }
        } else {
            if (this.state.allowMergeInit) {
                initMergeButton = <Nav.Link onClick={this.showMergeModal} className="icon-link"><IconMerge className="feather feather-lg" /></Nav.Link>;
            }
        }

        collapsable = (
            <>
                <Nav className="mr-auto">
                    <NavLink to="/search" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}>{this.props.t("bookSeat")}</NavLink>
                    <NavLink to="/bookings" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}>{this.props.t("myBookings")}</NavLink>
                    <NavLink to="/preferences" className={({isActive}) => "nav-link " + (isActive ? "active" : "")}>{this.props.t("preferences")}</NavLink>
                    {signOffButton}
                </Nav>
                <Nav className="mr-right">
                    {initMergeButton}
                    {mergeRequestsButton}
                    {userInfo}
                </Nav>
            </>
        );

        if (!RuntimeConfig.EMBEDDED) {
            collapsable = (
                <>
                    <Navbar.Toggle aria-controls="basic-navbar-nav" />
                    <Navbar.Collapse id="basic-navbar-nav">
                        {collapsable}
                    </Navbar.Collapse>
                </>
            );
        }

        return (
            <>
                <Navbar bg="dark" variant="dark" expand="lg" fixed="top" className={navClass}>
                    <Navbar.Brand as={NavLink} to="/search"><img src="/ui/seatsurfing_white.svg" alt="Seatsurfing" /></Navbar.Brand>
                    {collapsable}
                </Navbar>
                <Modal show={this.state.showMergeInit} onHide={() => this.setState({ showMergeInit: false })}>
                    <Modal.Header closeButton>
                        <Modal.Title>{this.props.t("mergeUserAccounts")}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <p>{this.props.t("mergeUserAccountsHint")}</p>
                        <Form onSubmit={(e) => { e.preventDefault(); this.initMerge() }}>
                            <Form.Group>
                                <Form.Control type="email" placeholder={this.props.t("emailPlaceholder")} required={true} value={this.state.targetUserEmail} onChange={(e: any) => this.setState({ targetUserEmail: e.target.value, invalidTargetUserEmail: false })} isInvalid={this.state.invalidTargetUserEmail} autoFocus={true} />
                                <Form.Control.Feedback type="invalid">{this.props.t("errorInvalidEmail")}</Form.Control.Feedback>
                            </Form.Group>
                        </Form>
                    </Modal.Body>
                    <Modal.Footer>
                        <Button variant="secondary" onClick={() => this.setState({ showMergeInit: false })}>
                            {this.props.t("cancel")}
                        </Button>
                        <Button variant="primary" onClick={this.initMerge}>
                            {this.props.t("requestMerge")}
                        </Button>
                    </Modal.Footer>
                </Modal>
                <Modal show={this.state.showMergeNextStep} onHide={() => this.setState({ showMergeNextStep: false })}>
                    <Modal.Header closeButton>
                        <Modal.Title>{this.props.t("mergeUserAccounts")}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <p>{this.props.t("mergeUserAccountsNextStepHint")}</p>
                        <Button onClick={this.openWebUI}>{this.props.t("openWebUI")}</Button>
                    </Modal.Body>
                </Modal>
                <Modal show={this.state.showMergeRequests} onHide={() => this.setState({ showMergeRequests: false })}>
                    <Modal.Header closeButton>
                        <Modal.Title>{this.props.t("mergeUserAccounts")}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <p>{this.props.t("introIncomingMergeRequests")}</p>
                        {this.state.mergeRequests.map(item => this.renderMergeRequest(item))}
                    </Modal.Body>
                </Modal>
            </>
        );
    }
}

export default withTranslation()(NavBar as any);
