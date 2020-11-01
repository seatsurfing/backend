import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { RouteChildrenProps, Link } from 'react-router-dom';
import { Search } from 'flexspace-commons';
import { Card, ListGroup, Col, Row } from 'react-bootstrap';

interface State {
  loading: boolean
}

interface Props {
  keyword: string
}

export default class SearchResult extends React.Component<RouteChildrenProps<Props>, State> {
  data: Search;

  constructor(props: any) {
    super(props);
    this.data = new Search();
    this.state = {
      loading: true
    };
  }

  componentDidMount = () => {
    this.loadItems();
  }

  componentDidUpdate = (prevProps: RouteChildrenProps<Props>) => {
    if (this.props.match?.params.keyword !== prevProps.match?.params.keyword) {
      this.loadItems();
    }
  }

  loadItems = () => {
    Search.search(this.props.match ? this.props.match.params.keyword : "").then(res => {
      this.data = res;
      this.setState({ loading: false });
    });
  }

  escapeHTML = (s: string): string => {
    return s;
    /*
    return s
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
    */
  }

  renderUserResults = () => {
    let items = this.data.users.map(user => {
      let link = "/users/" + user.id;
      return (
        <ListGroup.Item key={user.id}><Link to={link}>{user.email}</Link></ListGroup.Item>
      );
    });
    if (items.length === 0) {
      items.push(<ListGroup.Item key="users-no-results">Keine Ergebnisse.</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>Benutzer</Card.Header>
          <ListGroup variant="flush">
            {items}
          </ListGroup>
        </Card>
      </Col>
    );
  }

  renderLocationResults = () => {
    let items = this.data.locations.map(location => {
      let link = "/locations/" + location.id;
      return (
        <ListGroup.Item key={location.id}><Link to={link}>{location.name}</Link></ListGroup.Item>
      );
    });
    if (items.length === 0) {
      items.push(<ListGroup.Item key="locations-no-results">Keine Ergebnisse.</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>Bereiche</Card.Header>
          <ListGroup variant="flush">
            {items}
          </ListGroup>
        </Card>
      </Col>
    );
  }

  renderSpaceResults = () => {
    let items = this.data.spaces.map(space => {
      let link = "/locations/" + space.locationId;
      return (
        <ListGroup.Item key={space.id}><Link to={link}>{space.name}</Link></ListGroup.Item>
      );
    });
    if (items.length === 0) {
      items.push(<ListGroup.Item key="spaces-no-results">Keine Ergebnisse.</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>Plätze</Card.Header>
          <ListGroup variant="flush">
            {items}
          </ListGroup>
        </Card>
      </Col>
    );
  }

  render() {
    let headline = "Suche nach '" + this.escapeHTML(this.props.match ? this.props.match.params.keyword : "") + "'"

    if (this.state.loading) {
      return (
        <FullLayout headline={headline}>
          <Loading />
        </FullLayout>
      );
    }

    return (
      <FullLayout headline={headline}>
        <Row>
          {this.renderUserResults()}
          {this.renderLocationResults()}
          {this.renderSpaceResults()}
        </Row>
      </FullLayout>
    );
  }
}
