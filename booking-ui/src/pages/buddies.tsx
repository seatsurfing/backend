import React from 'react';
import { Ajax, Buddy } from 'flexspace-commons';
import Loading from '../components/Loading';
import { Button, Form, ListGroup, Modal } from 'react-bootstrap';
import { WithTranslation, withTranslation } from 'next-i18next';
import { NextRouter } from 'next/router';
import NavBar from '@/components/NavBar';
import withReadyRouter from '@/components/withReadyRouter';

interface State {
  loading: boolean
  selectedItem: Buddy | null
}

interface Props extends WithTranslation {
  router: NextRouter
}

class Buddies extends React.Component<Props, State> {
  data: Buddy[];

  constructor(props: any) {
    super(props);
    this.data = [];
    this.state = {
      loading: true,
      selectedItem: null
    };
  }

  componentDidMount = () => {
    if (!Ajax.CREDENTIALS.accessToken) {
      this.props.router.push("/login");
      return;
    }
    this.loadData();
  }

  loadData = () => {
    Buddy.list().then(list => {
      this.data = list;
      this.setState({ loading: false });
    });
  }

  onItemPress = (item: Buddy) => {
    this.setState({ selectedItem: item });
  }

  removeBuddy = (item: Buddy | null) => {
    this.setState({
      loading: true
    });
    this.state.selectedItem?.delete().then(() => {
      this.setState({
        selectedItem: null,
      }, this.loadData);
    });
  }

  renderItem = (item: Buddy) => {
    return (
      <ListGroup.Item key={item.id}>
        <h5>{item.buddy.email}</h5>
        <p>{this.props.t("nextBooking") + ": " + item.buddy.firstBooking}</p>
        <Button variant="danger" onClick={() => this.onItemPress(item)}>
          {this.props.t("removeBuddy")}
        </Button>
      </ListGroup.Item>
    );
  }

  render() {
    if (this.state.loading) {
      return <Loading />;
    }
    if (this.data.length === 0) {
      return (
        <>
          <NavBar />
          <div className="container-signin">
            <Form className="form-signin">
              <p>{this.props.t("noBuddies")}</p>
            </Form>
          </div>
        </>
      );
    }
    return (
      <>
        <NavBar />
        <div className="container-signin">
          <Form className="form-signin">
            <ListGroup>
              {this.data.map(item => this.renderItem(item))}
            </ListGroup>
          </Form>
        </div>
        <Modal show={this.state.selectedItem != null} onHide={() => this.setState({ selectedItem: null })}>
          <Modal.Header closeButton>
            <Modal.Title>{this.props.t("removeBuddy")}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>{this.props.t("confirmRemoveBuddy", { interpolation: { escapeValue: false } })}</p>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => this.setState({ selectedItem: null })}>
              {this.props.t("back")}
            </Button>
            <Button variant="danger" onClick={() => this.removeBuddy(this.state.selectedItem)}>
              {this.props.t("removeBuddy")}
            </Button>
          </Modal.Footer>
        </Modal>
      </>
    );
  }
}

export default withTranslation()(withReadyRouter(Buddies as any));
