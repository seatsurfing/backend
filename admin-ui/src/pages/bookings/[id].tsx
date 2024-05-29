import React from 'react';
import FullLayout from '../../components/FullLayout';
import { Form, Col, Row, Button, Alert, InputGroup } from 'react-bootstrap';
import { Navigate, NavigateFunction, Params, PathRouteProps } from 'react-router-dom';
import { ChevronLeft as IconBack, Save as IconSave, Trash2 as IconDelete } from 'react-feather';
import Loading from '../../components/Loading';
import { Location, Space, Booking, Formatting, User, AuthProvider, Settings as OrgSettings } from 'flexspace-commons';
import { NextRouter } from 'next/router';
import { WithTranslation, withTranslation } from 'next-i18next';
import { TFunction } from 'i18next';
import Link from 'next/link';
import withReadyRouter from '@/components/withReadyRouter';
import DateTimePicker from 'react-datetime-picker';
import DatePicker from 'react-date-picker';
import 'react-datetime-picker/dist/DateTimePicker.css';
import 'react-date-picker/dist/DatePicker.css';
import 'react-clock/dist/Clock.css';

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
    isDisabledLocation: boolean
    isDisabledSpace: boolean
    canSearch: boolean
    canSearchHint: string
}

interface Props extends WithTranslation {
    router: NextRouter
}

class EditBooking extends React.Component<Props, State> {
    entity: Booking = new Booking();
    authProviders: { [key: string]: string } = {};
    dailyBasisBooking: boolean;
    noAdminRestrictions: boolean;
    maxBookingsPerUser: number
    maxDaysInAdvance: number
    maxBookingDurationHours: number
    isNewBooking: boolean;
    enterChangeTimer: number | undefined;
    leaveChangeTimer: number | undefined;
    curBookingCount: number = 0;

    constructor(props: any) {
        super(props);
        this.dailyBasisBooking = false;
        this.noAdminRestrictions = false;
        this.maxBookingsPerUser = 0;
        this.maxBookingDurationHours = 0;
        this.maxDaysInAdvance = 0;
        this.isNewBooking = false;
        this.enterChangeTimer = undefined;
        this.leaveChangeTimer = undefined;
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
            isDisabledLocation: true,
            isDisabledSpace: true,
            canSearch: false,
            canSearchHint: "",
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

    loadData = async (): Promise<void> => {
        const {id} = this.props.router.query;
        if (id && (typeof id === "string")){
            if (id !== 'add') {
            return Booking.get(id).then(booking => {
                this.entity = booking;
                this.setState({
                    enter: this.entity.enter,
                    leave: this.entity.leave,
                    selectedLocationId: this.entity.space.locationId,
                    selectedSpaceId: this.entity.space.id,
                    selectedUserEmail: this.entity.user.email,
                    isDisabledLocation: false,
                    isDisabledSpace: false
                    // loading: false,
                });
                this.loadSpaces(this.entity.space.locationId, this.entity.enter, this.entity.leave);
                this.isNewBooking = false;
            });
            } else {
                // add 
                this.isNewBooking = true;
                this.setState({
                    isDisabledLocation: false,
                    isDisabledSpace: false
                    // loading: false,
                });

            }
        } else {
            // no id
        }
    }

    loadSpaces = async (selectedLocationId: string, enter: Date, leave: Date): Promise<void> => {
        // this.setState({ loading: true });
        return Space.listAvailability(selectedLocationId, enter, leave).then(list => {
            this.setState({ 
                spaces: list, 
                // loading: false
            });
        });
    }

    loadSettings = async (): Promise<void> => {
        return OrgSettings.list().then(settings => {
            settings.forEach(s => {
                if (s.name === "daily_basis_booking") {this.dailyBasisBooking = (s.value === "1")};
                if (s.name === "no_admin_restrictions") { this.noAdminRestrictions = (s.value === "1")};
                if (s.name === "max_bookings_per_user") {this.maxBookingsPerUser = window.parseInt(s.value)};
                if (s.name === "max_days_in_advance") {this.maxDaysInAdvance = window.parseInt(s.value)};
                if (s.name === "max_booking_duration_hours") {this.maxBookingDurationHours = window.parseInt(s.value)};
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

    //TODO: modify to init according to selcted user
    // initCurrentBookingCount = () => {
    //     Booking.list().then(list => {
    //         this.curBookingCount = list.length;
    //         this.updateCanSearch();
    //     });
    // }
    
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
            let booking = new Booking();
            booking.enter = this.state.enter;
            booking.leave = this.state.leave;
            booking.space.id = this.state.selectedSpaceId;
            booking.user.email = this.state.selectedUserEmail;
            booking.save().then(() => {
                // this.props.navigate("/bookings/");
                this.setState({ saved: true });
            }).catch(() => {
                this.setState({ 
                    error: true,
                    saved: false
                });
            });    
        } else {
            this.entity.enter = this.state.enter;
            this.entity.leave = this.state.leave;
            this.entity.space.id = this.state.selectedSpaceId;
            this.entity.user.email = this.state.selectedUserEmail;
            this.entity.save().then(() => {
                // this.props.navigate("/bookings/" + this.entity.id);
                this.loadData();
                this.setState({ saved: true });
            }).catch(() => {
                this.setState({
                    error: true,
                    saved: false
                });
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

    updateCanSearch = async () => {
        let res = true;
        let hint = "";
        if (this.curBookingCount >= this.maxBookingsPerUser) {
            res = false;
            hint = this.props.t("errorBookingLimit", { "num": this.maxBookingsPerUser });
        }
        // if (!this.state.selectedLocationId && !this.entity.location.id) {
        //     res = false;
        //     hint = this.props.t("errorPickArea");
        // }
        let now = new Date();
        let enterTime = new Date(this.state.enter);
        if (this.dailyBasisBooking) {
            enterTime.setHours(23, 59, 59);
        }
        if (enterTime.getTime() < now.getTime()) {
            res = false;
            hint = this.props.t("errorEnterFuture");
        }
        if (this.state.leave.getTime() <= this.state.enter.getTime()) {
            res = false;
            hint = this.props.t("errorLeaveAfterEnter");
        }
        const MS_PER_MINUTE = 1000 * 60;
        const MS_PER_HOUR = MS_PER_MINUTE * 60;
        const MS_PER_DAY = MS_PER_HOUR * 24;
        let bookingAdvanceDays = Math.floor((this.state.enter.getTime() - new Date().getTime()) / MS_PER_DAY);
        if (bookingAdvanceDays > this.maxDaysInAdvance && !this.noAdminRestrictions) {
            res = false;
            hint = this.props.t("errorDaysAdvance", { "num": this.maxDaysInAdvance });
        }
        let bookingDurationHours = Math.floor((this.state.leave.getTime() - this.state.enter.getTime()) / MS_PER_MINUTE) / 60;
        if (bookingDurationHours > this.maxBookingDurationHours && !this.noAdminRestrictions) {
            res = false;
            hint = this.props.t("errorBookingDuration", { "num": this.maxBookingDurationHours });
        }
        let self = this;
        return new Promise<void>(function (resolve, reject) {
            self.setState({
                canSearch: res,
                canSearchHint: hint
            }, () => resolve());
        });
    }

    setEnterDate = (value: Date | [Date | null, Date | null]) => {
        let dateChangedCb = () => {
            this.updateCanSearch().then(() => {
                if (!this.state.canSearch) {
                    this.setState({ loading: false });
                } else {
                    // let promises = [
                    //     this.initCurrentBookingCount(),
                    //     this.loadSpaces(this.state.locationId),
                    // ];
                    // Promise.all(promises).then(() => {
                    //     this.setState({ loading: false });
                    // });
                }
            });
        };
        let performChange = () => {
            let enter = (value instanceof Date) ? value : value[0];
            if (enter == null) {
              return;
            }
            let leave = new Date(enter);
            leave.setHours(leave.getHours()+1);
            if (this.dailyBasisBooking) {
                enter.setHours(0, 0, 0);
                leave.setHours(23, 59, 59);
            }
            this.setState({
                enter: enter,
                leave: leave,
                isDisabledLocation: false
            }, () => dateChangedCb());

            if (this.state.selectedLocationId) {
                this.loadSpaces(this.state.selectedLocationId, this.state.enter, this.state.leave)
            }
        };
        window.clearTimeout(this.leaveChangeTimer);
        this.leaveChangeTimer = window.setTimeout(performChange, 1000);
        return true;
    }

    setLeaveDate = (value: Date | [Date | null, Date | null]) => {
        let dateChangedCb = () => {
            //TODO: check for parameters *maxBookingDur ...

            this.updateCanSearch().then(() => {
                if (!this.state.canSearch) {
                    this.setState({ loading: false });
                } else {
                    // let promises = [
                    //     this.initCurrentBookingCount(),
                    //     this.loadSpaces(this.state.locationId),
                    // ];
                    // Promise.all(promises).then(() => {
                    //     this.setState({ loading: false });
                    // });
                }
            });
        };
        let performChange = () => {
            let date = (value instanceof Date) ? value : value[0];
            if (date == null) {
              return;
            }
            if (this.dailyBasisBooking) {
                date.setHours(23, 59, 59);
            }
            this.setState({
                leave: date,
                isDisabledLocation: false
            }, () => dateChangedCb());
            if (this.state.selectedLocationId) {
                this.loadSpaces(this.state.selectedLocationId, this.state.enter, date)
            }
        };
        window.clearTimeout(this.leaveChangeTimer);
        this.leaveChangeTimer = window.setTimeout(performChange, 1000);
    }
    
    render() {
        if (this.state.goBack) {
            return <Navigate replace={true} to={`/bookings`} />
        }

        let hint = <></>;
        if ((!this.state.canSearch) && (this.state.canSearchHint)) {
            hint = (
                <Form.Group as={Row} className="margin-top-10">
                <Col xs="2"></Col>
                <Col xs="10">
                    <div className="invalid-search-config">{this.state.canSearchHint}</div>
                </Col>
                </Form.Group>
            );
        }

        let enterDatePicker = <DateTimePicker value={this.state.enter} onChange={(value: Date | null) => { if (value != null) this.setEnterDate(value) }} clearIcon={null} required={true} format={this.props.t("datePickerFormat")} />;
        if (this.dailyBasisBooking) {
          enterDatePicker = <DatePicker value={this.state.enter} onChange={(value: Date | null | [Date | null, Date | null]) => { if (value != null) this.setEnterDate(value) }} clearIcon={null} required={true} format={this.props.t("datePickerFormatDailyBasisBooking")} />;
        }
        let leaveDatePicker = <DateTimePicker value={this.state.leave} onChange={(value: Date | null) => { if (value != null) this.setLeaveDate(value) }} clearIcon={null} required={true} format={this.props.t("datePickerFormat")} />;
        if (this.dailyBasisBooking) {
          leaveDatePicker = <DatePicker value={this.state.leave} onChange={(value: Date | null | [Date | null, Date | null]) => { if (value != null) this.setLeaveDate(value) }} clearIcon={null} required={true} format={this.props.t("datePickerFormatDailyBasisBooking")} />;
        }

        let backButton = <Link href="/bookings" className="btn btn-sm btn-outline-secondary"><IconBack className="feather" /> {this.props.t("back")}</Link>;
        let buttons = backButton;

        if (this.state.loading) {
            return (
                <FullLayout headline={this.props.t("editBooking")} buttons={buttons}>
                    <Loading />
                </FullLayout>
            );
        }

        if (this.state.saved) {
            hint = <Alert variant="success">{this.props.t("entryUpdated")}</Alert>
        } else if (this.state.canSearchHint) {
            hint = <Alert variant="danger">{this.props.t(this.state.canSearchHint)}</Alert>
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
            <FullLayout headline={this.props.t("editBooking")} buttons={buttons}>
                <Form onSubmit={this.onSubmit} id="form">

                    {hint}

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("user")}</Form.Label>
                        <Col sm="4">
                            <Form.Select required={true} value={this.state.selectedUserEmail} onChange={(e: any) => this.setState({ selectedUserEmail: e.target.value })}>
                                {/* TODO: if (this.entity.user.email) { */}
                                <option disabled={true} value={this.entity.user.id}>{this.entity.user.email}</option>                          
                                {this.state.users.map((user: {email: string | undefined; }) => (
                                    <option key={user.email} value={user.email}>{user.email}</option>
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
                            <Form.Select disabled={this.state.isDisabledLocation} required={true} value={this.state.selectedLocationId} onChange={(e: any) => {this.setState({ selectedLocationId: e.target.value, isDisabledSpace: false }); this.loadSpaces(e.target.value, this.state.enter, this.state.leave)}}>
                                <option disabled={true} value={this.entity.space.location.id}>{this.entity.space.location.name}</option>
                                {this.state.locations.map((location: {name: string | undefined; id: string | undefined}) => (
                                    <option key={location.id} value={location.id}>{location.name}</option>
                                ))}
                            </Form.Select>
                        </Col>
                    </Form.Group>

                    <Form.Group as={Row}>
                        <Form.Label column sm="2">{this.props.t("space")}</Form.Label>
                        <Col sm="4">
                            <Form.Select disabled={this.state.isDisabledSpace} required={true} value={this.state.selectedSpaceId} onChange={(e: any) => this.setState({ selectedSpaceId: e.target.value })}>
                                <option disabled={true} value={this.entity.space.id}>{this.entity.space.name}</option>
                                {this.state.spaces.map(function(space: { id: string | undefined; name: string | null | undefined; available: boolean}){ 
                                    if(space.available){
                                        return <option key={space.id} value={space.id}>{space.name}</option>
                                    }else{
                                        return <option key={space.id} disabled value={space.id}>{space.name}</option>
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

export default withTranslation(['admin'])(withReadyRouter(EditBooking as any));
