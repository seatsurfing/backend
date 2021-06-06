import React from 'react';
import './Login.css';
import { Form, Button, Alert } from 'react-bootstrap';
import { Location, Booking } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
// @ts-ignore
import DateTimePicker from 'react-datetime-picker';
import DatePicker from 'react-date-picker';
import './Search.css';
import { Redirect } from 'react-router-dom';
import { SearchResultRouteParams } from './SearchResult';
import { AuthContext } from '../AuthContextData';

interface State {
  enter: Date
  leave: Date
  locations: Location[]
  locationId: string
  canSearch: boolean
  canSearchHint: string
  showResult: boolean
}

interface Props {
  t: TFunction
}
class Search extends React.Component<Props, State> {
  static contextType = AuthContext;
  curBookingCount: number = 0;

  constructor(props: any) {
    super(props);
    this.state = {
      enter: new Date(),
      leave: new Date(),
      locations: [],
      locationId: "",
      canSearch: false,
      canSearchHint: "",
      showResult: false
    };
  }

  componentDidMount = () => {
    this.initDates();
    this.initCurrentBookingCount();
    Location.list().then(list => {
      this.setState({ locations: list });
    });
  }

  initCurrentBookingCount = () => {
    Booking.list().then(list => {
      this.curBookingCount = list.length;
      this.updateCanSearch();
    });
  }

  initDates = () => {
    let now = new Date();
    if (now.getHours() > 17) {
      let enter = new Date();
      enter.setDate(enter.getDate() + 1);
      if (this.context.dailyBasisBooking) {
        enter.setHours(0, 0, 0);
      } else {
        enter.setHours(9, 0, 0);
      }
      let leave = new Date(enter);
      if (this.context.dailyBasisBooking) {
        leave.setHours(23, 59, 59);
      } else {
        leave.setHours(17, 0, 0);
      }
      this.setState({
        enter: enter,
        leave: leave
      });
    } else {
      if (this.context.dailyBasisBooking) {
        let enter = new Date();
        enter.setHours(0, 0, 0);
        let leave = new Date(enter);
        leave.setHours(23, 59, 59);
        this.setState({
          enter: enter,
          leave: leave
        });
      } else {
        let enter = new Date();
        enter.setHours(enter.getHours() + 1, 0, 0);
        let leave = new Date(enter);
        if (leave.getHours() < 17) {
          leave.setHours(17, 0, 0);
        } else {
          leave.setHours(leave.getHours() + 1, 0, 0);
        }
        this.setState({
          enter: enter,
          leave: leave
        });
      }
    }
  }

  updateCanSearch = () => {
    let res = true;
    let hint = "";
    if (this.curBookingCount >= this.context.maxBookingsPerUser) {
      res = false;
      hint = this.props.t("errorBookingLimit", { "num": this.context.maxBookingsPerUser });
    }
    if (!this.state.locationId) {
      res = false;
      hint = this.props.t("errorPickArea");
    }
    let now = new Date();
    let enterTime = new Date(this.state.enter);
    if (this.context.dailyBasisBooking) {
      enterTime.setHours(23, 59, 59);
    }
    if (enterTime.getTime() <= now.getTime()) {
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
    if (bookingAdvanceDays > this.context.maxDaysInAdvance) {
      res = false;
      hint = this.props.t("errorDaysAdvance", { "num": this.context.maxDaysInAdvance });
    }
    let bookingDurationHours = Math.floor((this.state.leave.getTime() - this.state.enter.getTime()) / MS_PER_MINUTE) / 60;
    if (bookingDurationHours > this.context.maxBookingDurationHours) {
      res = false;
      hint = this.props.t("errorBookingDuration", { "num": this.context.maxBookingDurationHours });
    }
    this.setState({
      canSearch: res,
      canSearchHint: hint
    });
  }

  renderLocations = () => {
    return this.state.locations.map(location => {
      return <option value={location.id} key={location.id}>{location.name}</option>;
    });
  }

  setEnterDate = (value: Date | Date[]) => {
    let date = (value instanceof Date) ? value : value[0];
    if (this.context.dailyBasisBooking) {
      date.setHours(0, 0, 0);
    }
    this.setState({ enter: date }, this.updateCanSearch);
  }

  setLeaveDate = (value: Date | Date[]) => {
    let date = (value instanceof Date) ? value : value[0];
    if (this.context.dailyBasisBooking) {
      date.setHours(23, 59, 59);
    }
    this.setState({ leave: date }, this.updateCanSearch);
  }

  setLocationId = (value: string) => {
    this.setState({ locationId: value }, this.updateCanSearch);
  }

  onSubmit = () => {
    this.setState({showResult: true});
  }

  render() {
    if (this.state.showResult) {
      let props: SearchResultRouteParams = {
        locationId: this.state.locationId,
        enter: this.state.enter,
        leave: this.state.leave
      };
      return <Redirect to={{pathname: "/search/result", state: props}} />
    }

    let hint = <></>;
    if (!this.state.canSearch) {
      hint = (
        <Form.Group>
          <Alert variant="warning">{this.state.canSearchHint}</Alert>
        </Form.Group>
      );
    } else {
      hint = (
        <Form.Group>
          <Button variant="primary" type="submit" disabled={!this.state.canSearch}>{this.props.t("searchSpace")}</Button>
        </Form.Group>
      );
    }
    let enterDatePicker = <DateTimePicker value={this.state.enter} onChange={(value: Date) => this.setEnterDate(value)} clearIcon={null} required={true} />;
    if (this.context.dailyBasisBooking) {
      enterDatePicker = <DatePicker value={this.state.enter} onChange={(value: Date | Date[]) => this.setEnterDate(value)} clearIcon={null} required={true} />;
    }
    let leaveDatePicker = <DateTimePicker value={this.state.leave} onChange={(value: Date) => this.setLeaveDate(value)} clearIcon={null} required={true} />;
    if (this.context.dailyBasisBooking) {
      leaveDatePicker = <DatePicker value={this.state.leave} onChange={(value: Date | Date[]) => this.setLeaveDate(value)} clearIcon={null} required={true} />;
    }
    return (
      <div className="container-signin">
        <Form className="form-signin" onSubmit={this.onSubmit}>
          <Form.Group>
            <Form.Label>{this.props.t("enter")}</Form.Label>
            {enterDatePicker}
          </Form.Group>
          <Form.Group>
            <Form.Label>{this.props.t("leave")}</Form.Label>
            {leaveDatePicker}
          </Form.Group>
          <Form.Group>
            <Form.Label>{this.props.t("area")}</Form.Label>
            <Form.Control as="select" custom={true} required={true} onChange={(e) => this.setLocationId(e.target.value)}>
              <option value="">({this.props.t("pleaseSelect")})</option>
              {this.renderLocations()}
            </Form.Control>
          </Form.Group>
          {hint}
        </Form>
      </div>
    )
  }
}

export default withTranslation()(Search as any);
