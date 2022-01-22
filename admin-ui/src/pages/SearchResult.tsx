import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { Link, Params, RouteProps } from 'react-router-dom';
import { Search } from 'flexspace-commons';
import { Card, ListGroup, Col, Row } from 'react-bootstrap';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { withRouter } from '../types/withRouter';

interface State {
  loading: boolean
}

interface Props extends RouteProps {
  params: Readonly<Params<string>>
  t: TFunction
}

class SearchResult extends React.Component<Props, State> {
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

  componentDidUpdate = (prevProps: Props) => {
    if (this.props.params.keyword !== prevProps.params.keyword) {
      this.loadItems();
    }
  }

  loadItems = () => {
    Search.search(this.props.params.keyword ? this.props.params.keyword : "").then(res => {
      this.data = res;
      this.setState({ loading: false });
    });
  }

  escapeHTML = (s: string): string => {
    return s;
  }

  renderUserResults = () => {
    let items = this.data.users.map(user => {
      let link = "/users/" + user.id;
      return (
        <ListGroup.Item key={user.id}><Link to={link}>{user.email}</Link></ListGroup.Item>
      );
    });
    if (items.length === 0) {
      items.push(<ListGroup.Item key="users-no-results">{this.props.t("noResults")}</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>{this.props.t("users")}</Card.Header>
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
      items.push(<ListGroup.Item key="locations-no-results">{this.props.t("noResults")}</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>{this.props.t("areas")}</Card.Header>
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
      items.push(<ListGroup.Item key="spaces-no-results">{this.props.t("noResults")}</ListGroup.Item>);
    }
    return (
      <Col sm="4" className="mb-4">
        <Card>
          <Card.Header>{this.props.t("spaces")}</Card.Header>
          <ListGroup variant="flush">
            {items}
          </ListGroup>
        </Card>
      </Col>
    );
  }

  render() {
    let headline = this.props.t("searchForX", {keyword: this.escapeHTML(this.props.params.keyword ? this.props.params.keyword : "")});

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

export default withRouter(withTranslation()(SearchResult as any));
