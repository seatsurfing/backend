import React from 'react';
import { Ajax, Location, UserPreference } from 'flexspace-commons';
import Loading from '../components/Loading';
import { Alert, Button, Col, Form, Row } from 'react-bootstrap';
import { WithTranslation, withTranslation } from 'next-i18next';
import { NextRouter } from 'next/router';
import NavBar from '@/components/NavBar';
import withReadyRouter from '@/components/withReadyRouter';

interface State {
  loading: boolean
  submitting: boolean
  saved: boolean
  error: boolean
  enterTime: number
  workdayStart: number
  workdayEnd: number
  workdays: boolean[]
  locationId: string
}

interface Props extends WithTranslation {
  router: NextRouter
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
      workdayStart: 0,
      workdayEnd: 0,
      workdays: [],
      locationId: "",
    };
  }

  componentDidMount = () => {
    if (!Ajax.CREDENTIALS.accessToken) {
      this.props.router.push("/login");
      return;
    }
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
          if (typeof window !== 'undefined') {
            if (s.name === "enter_time") state.enterTime = window.parseInt(s.value);
            if (s.name === "workday_start") state.workdayStart = window.parseInt(s.value);
            if (s.name === "workday_end") state.workdayEnd = window.parseInt(s.value);
          }
          if (s.name === "workdays") {
            state.workdays = [];
            for (let i = 0; i <= 6; i++) {
              state.workdays[i] = false;
            }
            s.value.split(",").forEach(val => state.workdays[val] = true)
          }
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
    let workdays: string[] = [];
    this.state.workdays.forEach((val, day) => {
      if (val) {
        workdays.push(day.toString());
      }
    });
    let payload = [
      new UserPreference("enter_time", this.state.enterTime.toString()),
      new UserPreference("workday_start", this.state.workdayStart.toString()),
      new UserPreference("workday_end", this.state.workdayEnd.toString()),
      new UserPreference("workdays", workdays.join(",")),
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

  onWorkdayCheck = (day: number, checked: boolean) => {
    let workdays = this.state.workdays.map((val, i) => (i === day) ? checked : val);
    this.setState({
      workdays: workdays
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
        <NavBar />
        <div className="container-center">
          <Form className="container-center-inner" onSubmit={this.onSubmit}>
            {hint}
            <Form.Group className="margin-top-15">
              <Form.Label>{this.props.t("notice")}</Form.Label>
              <Form.Select value={this.state.enterTime} onChange={(e: any) => this.setState({ enterTime: e.target.value })}>
                <option value="1">{this.props.t("earliestPossible")}</option>
                <option value="2">{this.props.t("nextDay")}</option>
                <option value="3">{this.props.t("nextWorkday")}</option>
              </Form.Select>
            </Form.Group>
            <Form.Group className="margin-top-15">
              <Form.Label>{this.props.t("workingHours")}</Form.Label>
              <Row>
                <Col>
                  <Form.Control type="number" value={this.state.workdayStart} onChange={(e: any) => this.setState({ workdayStart: typeof window !== 'undefined' ? window.parseInt(e.target.value) : 0 })} min="0" max="23" />
                </Col>
                <Col>
                  <Form.Control plaintext={true} readOnly={true} defaultValue={this.props.t("to").toString()} />
                </Col>
                <Col>
                  <Form.Control type="number" value={this.state.workdayEnd} onChange={(e: any) => this.setState({ workdayEnd: e.target.value })} min={this.state.workdayStart+1} max="23" />
                </Col>
              </Row>
            </Form.Group>
            <Form.Group className="margin-top-15">
              <Form.Label>{this.props.t("workdays")}</Form.Label>
              <div className="text-left">
                {[0, 1, 2, 3, 4, 5, 6].map(day => (
                  <Form.Check type="checkbox" key={"workday-" + day} id={"workday-" + day} label={this.props.t("workday-" + day)} checked={this.state.workdays[day]} onChange={(e: any) => this.onWorkdayCheck(day, e.target.checked)} />
                ))}
              </div>
            </Form.Group>
            <Form.Group className="margin-top-15">
              <Form.Label>{this.props.t("preferredLocation")}</Form.Label>
              <Form.Select value={this.state.locationId} onChange={(e: any) => this.setState({ locationId: e.target.value })}>
                <option value="">({this.props.t("none")})</option>
                {this.locations.map(location => <option key={"location-" + location.id} value={location.id}>{location.name}</option>)}
              </Form.Select>
            </Form.Group>
            <Button className="margin-top-15" type="submit">{this.props.t("save")}</Button>
          </Form>
        </div>
      </>
    );
  }
}

export default withTranslation()(withReadyRouter(Preferences as any));
