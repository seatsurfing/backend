import React from 'react';
import { Location, UserPreference } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { RouteProps } from 'react-router-dom';
import Loading from '../components/Loading';
import { Alert, Button, Form } from 'react-bootstrap';

interface State {
  loading: boolean
  submitting: boolean
  saved: boolean
  error: boolean
  enterTime: number
  bookingDuration: number
  locationId: string
}

interface Props extends RouteProps {
  t: TFunction
}

class Preferences extends React.Component<Props, State> {
  locations: Location[];

  constructor(props: any) {
    super(props);
    this.locations = [];
    this.state = {
      loading: true,
      submitting: false,
      saved: false,
      error: false,
      enterTime: 0,
      bookingDuration: 0,
      locationId: "",
    };
  }

  componentDidMount = () => {
    let promises = [
      this.loadPreferences(),
      this.loadLocations(),
    ];
    Promise.all(promises).then(() => {
      this.setState({ loading: false });
    });
  }

  loadPreferences = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      UserPreference.list().then(list => {
        let state: any = {};
        list.forEach(s => {
          if (s.name === "enter_time") state.enterTime = window.parseInt(s.value);
          if (s.name === "booking_duration") state.bookingDuration = window.parseInt(s.value);
          if (s.name === "location_id") state.locationId = s.value;
        });
        self.setState({
          ...self.state,
          ...state
        }, () => resolve());
      }).catch(e => reject(e));
    });
  }

  loadLocations = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      Location.list().then(list => {
        self.locations = list;
        resolve();
      }).catch(e => reject(e));
    });
  }

  onSubmit = (e: any) => {
    e.preventDefault();
    this.setState({
      submitting: true,
      saved: false,
      error: false
    });
    let payload = [
      new UserPreference("enter_time", this.state.enterTime.toString()),
      new UserPreference("booking_duration", this.state.bookingDuration.toString()),
      new UserPreference("location_id", this.state.locationId),
    ];
    UserPreference.setAll(payload).then(() => {
      this.setState({
        submitting: false,
        saved: true
      });
    }).catch(() => {
      this.setState({
        submitting: false,
        error: true
      });
    });
  }

  render() {
    if (this.state.loading) {
      return <Loading />;
    }

    let hint = <></>;
    if (this.state.saved) {
      hint = <Alert variant="success">{this.props.t("entryUpdated")}</Alert>
    } else if (this.state.error) {
      hint = <Alert variant="danger">{this.props.t("errorSave")}</Alert>
    }

    return (
      <>
        <div className="container-center">
          <Form className="container-center-inner" onSubmit={this.onSubmit}>
            {hint}
            <Form.Group>
              <Form.Label>{this.props.t("notice")}</Form.Label>
              <Form.Control as="select" custom={true} value={this.state.enterTime} onChange={(e: any) => this.setState({ enterTime: e.target.value })}>
                <option value="1">{this.props.t("earliestPossible")}</option>
                <option value="2">{this.props.t("nextDay")}</option>
                <option value="3">{this.props.t("nextWorkday")}</option>
              </Form.Control>
            </Form.Group>
            <Form.Group>
              <Form.Label>{this.props.t("bookingDuration")}</Form.Label>
              <Form.Control type="number" value={this.state.bookingDuration} onChange={(e: any) => this.setState({ bookingDuration: e.target.value })} min="1" max="9999" />
            </Form.Group>
            <Form.Group>
              <Form.Label>{this.props.t("preferredLocation")}</Form.Label>
              <Form.Control as="select" custom={true} value={this.state.locationId} onChange={(e: any) => this.setState({ locationId: e.target.value })}>
                <option value="">({this.props.t("none")})</option>
                {this.locations.map(location => <option value={location.id}>{location.name}</option>)}
              </Form.Control>
            </Form.Group>
            <Button type="submit">{this.props.t("save")}</Button>
          </Form>
        </div>
      </>
    );
  }
}

export default withTranslation()(Preferences as any);
