import React from 'react';
import FullLayout from '../components/FullLayout';
import { Form, Col, Row, Button, Alert, InputGroup } from 'react-bootstrap';
import { Link, Navigate, NavigateFunction, Params, PathRouteProps } from 'react-router-dom';
import { ChevronLeft as IconBack, Save as IconSave, Trash2 as IconDelete } from 'react-feather';
import Loading from '../components/Loading';
import { Location, Space, Booking, Formatting, User, AuthProvider, Settings as OrgSettings } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { withRouter } from '../types/withRouter';
import { withNavigate } from '../types/withNavigate';
// @ts-ignore
import DateTimePicker from 'react-datetime-picker';
import DatePicker from 'react-date-picker';
import './EditBooking.css';

interface State {
    loading: boolean
    submitting: boolean
    saved: boolean
    error: boolean
    goBack: boolean
    enter: Date
    leave: Date
    location: Location
    space: Space
    user: User
    selectedUserEmail: string
    selectedLocationId: string
    selectedSpaceId: string
    users: User[]
    locations: Location[]
    spaces: Space[]
}

interface Props extends PathRouteProps {
    navigate: NavigateFunction
    params: Readonly<Params<string>>
    t: TFunction
}

class EditBooking extends React.Component<Props, State> {
    entity: Booking = new Booking();
    authProviders: { [key: string]: string } = {};
    dailyBasisBooking: boolean
    isNewBooking: boolean

    constructor(props: any) {
        super(props);
        this.dailyBasisBooking = false;
        this.isNewBooking = false;
        this.state = {
            loading: true,
            submitting: false,
            saved: false,
            error: false,
            goBack: false,
            enter: new Date(),
            leave: new Date(),
            location: new Location(),
            space: new Space(),
            user: new User(),
            selectedUserEmail: "",
            selectedLocationId: "",
            selectedSpaceId: "",
            users: [],
            locations: [],
            spaces: [],
        }
    }

    componentDidMount = () => {
        let promises = [
            this.loadData(),
            this.loadSettings(),
            this.loadUsers(),
            this.loadLocations()
          ];
          Promise.all(promises).then(() => {
            this.setState({ loading: false });
          });
    }

    // TODO: load locations maxBookings !!!

    loadData = async (id?: string): Promise<void> => {
        if (!id) {
            id = this.props.params.id;
            this.isNewBooking = true;
        }
        if (id) {
            return Booking.get(id).then(booking => {
                this.entity = booking;
                this.setState({
                    enter: this.entity.enter,
                    leave: this.entity.leave,
                    selectedLocationId: this.entity.space.locationId,
                    selectedSpaceId: this.entity.space.id,
                    selectedUserEmail: this.entity.user.email,
                    loading: false,
                });
                this.loadSpaces(this.entity.space.locationId, this.entity.enter, this.entity.leave);
                this.isNewBooking = false;
            });
        } else {
            //return Promise.resolve();
        }
    }

    loadSpaces = async (selectedLocationId: string, enter: Date, leave: Date): Promise<void> => {
        this.setState({ loading: true });
        return Space.listAvailability(selectedLocationId, enter, leave).then(list => {
            this.setState({ 
                spaces: list, 
                loading: false });
        });
    }

    loadSettings = async (): Promise<void> => {
        return OrgSettings.list().then(settings => {
            settings.forEach(s => {
                if (s.name === "daily_basis_booking") {
                    this.dailyBasisBooking = (s.value === "1")
                }
                // this.setState({ loading: false });
            });
        });
    }
 
    loadUsers = () => {
        AuthProvider.list().then(providers => {
            providers.forEach(provider => {
                this.authProviders[provider.id] = provider.name;
            });
            User.list().then(list => {
                this.setState({ users: list })
                // this.setState({ loading: false });
            });
        });
    }

    loadLocations = async (): Promise<void> => {
        return Location.list().then(list => {
            this.setState({ locations: list })
            // this.setState({ loading: false });
        });
    }
    
    onSubmit = (e: any) => {
        e.preventDefault();
        this.setState({
            error: false,
            saved: false
        });

        if (this.dailyBasisBooking) {
            let enter = new Date();
            enter = this.state.enter;
            enter.setHours(0, 0, 0, 0)

            let leave = new Date();
            leave = this.state.leave;
            leave.setHours(23, 59, 59, 0)
  
            this.setState({
                enter: enter,
                leave: leave
            });
        }

        if (this.isNewBooking) {
            console.log("saving new booking: this.isNewBooking", this.isNewBooking);
            let booking = new Booking();
            booking.enter = this.state.enter;
            booking.leave = this.state.leave;
            booking.space.id = this.state.selectedSpaceId;
            booking.user.email = this.state.selectedUserEmail;
            booking.save().then(() => {
                // this.props.navigate("/bookings/");
                this.setState({ saved: true });
            }).catch(() => {
                this.setState({ error: true });
            });    
        } else {
            console.log("updating existing booking: this.isNewBooking", this.isNewBooking);
            this.entity.enter = this.state.enter;
            this.entity.leave = this.state.leave;
            this.entity.space.id = this.state.selectedSpaceId;
            this.entity.user.email = this.state.selectedUserEmail;
            console.log(this.entity);
            this.entity.update().then(() => {
                this.props.navigate("/bookings/" + this.entity.id);
                this.setState({ saved: true });
            }).catch(() => {
                this.setState({ error: true });
            });
        }
    }

    deleteItem = () => {
        if (window.confirm(this.props.t("confirmCancelBooking"))) {
            this.entity.delete().then(() => {
            this.setState({ goBack: true });
            });
        }
    }

    render() {
        if (this.state.goBack) {
            return <Navigate replace={true} to={`/bookings`} />
        }

        let enterDatePicker = <DateTimePicker value={this.state.enter} onChange={(value: Date) => this.setState({enter: value})} clearIcon={null} required={true}  />;
        if (this.dailyBasisBooking) {
            enterDatePicker = <DatePicker value={this.state.enter} onChange={(value: Date ) => this.setState({enter: value})} clearIcon={null} required={true}  />;
        }
        let leaveDatePicker = <DateTimePicker value={this.state.leave} onChange={(value: Date) => this.setState({leave: value})} clearIcon={null} required={true}  />;
        if (this.dailyBasisBooking) {
            leaveDatePicker = <DatePicker value={this.state.leave} onChange={(value: Date ) => this.setState({leave: value})} clearIcon={null} required={true}  />;
        }

        let backButton = <Link to="/bookings" className="btn btn-sm btn-outline-secondary"><IconBack className="feather" /> {this.props.t("back")}</Link>;
        let buttons = backButton;

        if (this.state.loading) {
            return (
                // TODO: add to TFunction
                <FullLayout headline={"Edit booking"} buttons={buttons}>
                    <Loading />
                </FullLayout>
            );
        }

        let hint = <></>;
        if (this.state.saved) {
            hint = <Alert variant="success">{this.props.t("entryUpdated")}</Alert>
        } else if (this.state.error) {
            hint = <Alert variant="danger">{this.props.t("errorSave")}</Alert>
        }

        let buttonDelete = <Button className="btn-sm" variant="outline-secondary" onClick={this.deleteItem} disabled={false}><IconDelete className="feather" /> {this.props.t("delete")}</Button>;
        let buttonSave = <Button className="btn-sm" variant="outline-secondary" type="submit" form="form"><IconSave className="feather" /> {this.props.t("save")}</Button>;
        if (this.entity.id) {
            buttons = <>{backButton} {buttonDelete} {buttonSave}</>;
        } else {
            buttons = <>{backButton} {buttonSave}</>;
        }

        return (
            // TODO: add to TFunction
            <FullLayout headline={"Edit booking"} buttons={buttons}>
                <Form onSubmit={this.onSubmit} id="form">
                    {hint}

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("user")}</Form.Label>
                        <Col sm="4">
                            <Form.Select required={true} value={this.state.selectedUserEmail} onChange={(e: any) => this.setState({ selectedUserEmail: e.target.value })}>
                                {/* TODO: if (this.entity.user.email) { */}
                                <option disabled={true} value={this.entity.user.id}>{this.entity.user.email}</option>                          
                                {this.state.users.map((user: {email: string | undefined; }) => (
                                    <option value={user.email}>{user.email}</option>
                                ))}
                            </Form.Select>
                        </Col>
                    </Form.Group>

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("enter")}</Form.Label>
                        <Col sm="4">
                            {enterDatePicker}
                        </Col>
                    </Form.Group>

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("leave")}</Form.Label>
                        <Col sm="4">
                            {leaveDatePicker}
                        </Col>
                    </Form.Group>

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("location")}</Form.Label>
                        <Col sm="4">
                            <Form.Select required={true} value={this.state.selectedLocationId} onChange={(e: any) => {this.setState({ selectedLocationId: e.target.value }); this.loadSpaces(e.target.value, this.state.enter, this.state.leave)}}>
                                <option disabled={true} value={this.entity.space.location.id}>{this.entity.space.location.name}</option>
                                {this.state.locations.map((location: {name: string | undefined; id: string | undefined}) => (
                                    <option value={location.id}>{location.name}</option>
                                ))}
                            </Form.Select>
                        </Col>
                    </Form.Group>

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("space")}</Form.Label>
                        <Col sm="4">
                            <Form.Select required={true} value={this.state.selectedSpaceId} onChange={(e: any) => this.setState({ selectedSpaceId: e.target.value })}>
                                <option disabled={true} value={this.entity.space.id}>{this.entity.space.name}</option>
                                {this.state.spaces.map(function(space: { id: string | undefined; name: string | null | undefined; available: boolean}){ 
                                    if(space.available){
                                        return <option value={space.id}>{space.name}</option>
                                    }else{
                                        return <option disabled value={space.id}>{space.name}</option>
                                    }
                                })}
                            </Form.Select>
                        </Col>
                    </Form.Group>

                </Form>
            </FullLayout>
        );

    }

}

export default withNavigate(withRouter(withTranslation()(EditBooking as any)));