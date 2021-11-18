import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { Stats, User } from 'flexspace-commons';
import { Card, Row, Col, ProgressBar } from 'react-bootstrap';
import './Dashboard.css';
import { Redirect } from 'react-router-dom';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
  loading: boolean
  redirect: string
  spaceAdmin: boolean
  orgAdmin: boolean
}

interface Props {
  t: TFunction
}

class Dashboard extends React.Component<Props, State> {
  stats: Stats | null;

  constructor(props: any) {
    super(props);
    this.stats = null;
    this.state = {
      loading: true,
      redirect: "",
      spaceAdmin: false,
      orgAdmin: false,
    };
  }

  componentDidMount = () => {
    let promises = [
      this.loadItems(),
      this.getUserInfo()
    ];
    Promise.all(promises).then(() => {
      this.setState({ loading: false });
    });
  }

  getUserInfo = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      User.getSelf().then(user => {
        self.setState({
          spaceAdmin: user.spaceAdmin,
          orgAdmin: user.admin,
        }, () => resolve());
      }).catch(e => reject(e));
    });
  }

  loadItems = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      Stats.get().then(stats => {
        self.stats = stats;
        resolve();
      }).catch(e => reject(e));
    });
  }

  renderStatsCard = (num: number | undefined, title: string, link?: string) => {
    let redirect = "";
    if (link) {
      redirect = link;
    }
    return (
      <Col sm="2">
        <Card className="dashboard-card-clickable" onClick={() => this.setState({ redirect: redirect })}>
          <Card.Body>
            <Card.Title className="dashboard-number text-center">{num}</Card.Title>
            <Card.Subtitle className="text-center mb-2 text-muted">{title}</Card.Subtitle>
          </Card.Body>
        </Card>
      </Col>
    );
  }

  renderProgressBar = (num: number | undefined, title: string) => {
    if (!num) {
      num = 0;
    }
    let label = title + ": " + num + " %";
    let variant = "success";
    if (num >= 80) {
      variant = "danger";
    }
    if (num >= 60) {
      variant = "warning";
    }
    return (
      <div>
        {label} <ProgressBar now={num} className="mb-3" variant={variant} />
      </div>
    );
  }

  render() {
    if (this.state.redirect) {
      return <Redirect to={this.state.redirect} />
    }

    if (this.state.loading) {
      return (
        <FullLayout headline="Dashboard">
          <Loading />
        </FullLayout>
      );
    }

    return (
      <FullLayout headline="Dashboard">
        <Row className="mb-4">
          {this.renderStatsCard(this.stats?.numUsers, this.props.t("users"), (this.state.orgAdmin ? "/users/": ""))}
          {this.renderStatsCard(this.stats?.numLocations, this.props.t("areas"), "/locations/")}
          {this.renderStatsCard(this.stats?.numSpaces, this.props.t("spaces"), "/locations/")}
          {this.renderStatsCard(this.stats?.numBookings, this.props.t("bookings"), "/bookings/")}
        </Row>
        <Row className="mb-4">
          {this.renderStatsCard(this.stats?.numBookingsToday, this.props.t("today"), "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsYesterday, this.props.t("yesterday"), "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsThisWeek, this.props.t("thisWeek"), "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsLastWeek, this.props.t("lastWeek"), "/bookings/")}
        </Row>
        <Row className="mb-4">
          <Col sm="8">
            <Card>
              <Card.Body>
                <Card.Title>{this.props.t("utilization")}</Card.Title>
                {this.renderProgressBar(this.stats?.spaceLoadToday, this.props.t("today"))}
                {this.renderProgressBar(this.stats?.spaceLoadYesterday, this.props.t("yesterday"))}
                {this.renderProgressBar(this.stats?.spaceLoadThisWeek, this.props.t("thisWeek"))}
                {this.renderProgressBar(this.stats?.spaceLoadLastWeek, this.props.t("lastWeek"))}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </FullLayout>
    );
  }
}

export default withTranslation()(Dashboard as any);
