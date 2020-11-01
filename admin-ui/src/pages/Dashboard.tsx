import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { Stats } from 'flexspace-commons';
import { Card, Row, Col, ProgressBar } from 'react-bootstrap';
import './Dashboard.css';
import { Redirect } from 'react-router-dom';

interface State {
  loading: boolean
  redirect: string
}

export default class Dashboard extends React.Component<{}, State> {
  stats: Stats | null;

  constructor(props: any) {
    super(props);
    this.stats = null;
    this.state = {
      loading: true,
      redirect: ""
    };
  }

  componentDidMount = () => {
    this.loadItems();
  }

  loadItems = () => {
    Stats.get().then(stats => {
      this.stats = stats;
      this.setState({ loading: false });
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
          {this.renderStatsCard(this.stats?.numUsers, "Benutzer", "/users/")}
          {this.renderStatsCard(this.stats?.numLocations, "Bereiche", "/locations/")}
          {this.renderStatsCard(this.stats?.numSpaces, "Plätze", "/locations/")}
          {this.renderStatsCard(this.stats?.numBookings, "Buchungen", "/bookings/")}
        </Row>
        <Row className="mb-4">
        {this.renderStatsCard(this.stats?.numBookingsToday, "Heute", "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsYesterday, "Gestern", "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsThisWeek, "Diese Woche", "/bookings/")}
          {this.renderStatsCard(this.stats?.numBookingsLastWeek, "Letzte Woche", "/bookings/")}
        </Row>
        <Row className="mb-4">
          <Col sm="8">
            <Card>
              <Card.Body>
                <Card.Title>Auslastung</Card.Title>
                  {this.renderProgressBar(this.stats?.spaceLoadToday, "Heute")}
                  {this.renderProgressBar(this.stats?.spaceLoadYesterday, "Gestern")}
                  {this.renderProgressBar(this.stats?.spaceLoadThisWeek, "Diese Woche")}
                  {this.renderProgressBar(this.stats?.spaceLoadLastWeek, "Letzte Woche")}
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </FullLayout>
    );
  }
}
