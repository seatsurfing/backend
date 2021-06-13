import React from 'react';
import './SearchResult.css';
import { Booking, Formatting, Location, Space } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
import { Link, RouteProps } from 'react-router-dom';
import Loading from '../components/Loading';
import { Button, Modal } from 'react-bootstrap';

interface State {
  loading: boolean
  locationId: string
  enter: Date
  leave: Date
  name: string
  showConfirm: boolean
  showSuccess: boolean
  showBookingNames: boolean
  selectedSpace: Space
}

export interface SearchResultRouteParams {
  locationId: string
  enter: Date
  leave: Date
}

interface Props extends RouteProps {
  t: TFunction
}

class SearchResult extends React.Component<Props, State> {
  data: Space[] = [];
  entity: Location = new Location();
  mapData: any = null;

  constructor(props: any) {
    super(props);
    this.state = {
      loading: true,
      locationId: "",
      enter: new Date(),
      leave: new Date(),
      name: "",
      showConfirm: false,
      showSuccess: false,
      showBookingNames: false,
      selectedSpace: new Space()
    };
  }

  componentDidMount = () => {
    if (this.props.location && this.props.location.state) {
      let data: SearchResultRouteParams = (this.props.location.state as SearchResultRouteParams);
      this.setState({
        locationId: data.locationId,
        enter: data.enter,
        leave: data.leave
      }, this.loadData);
    }
  }

  loadData = () => {
    Location.get(this.state.locationId).then(location => {
      this.entity = location;
      Space.listAvailability(location.id, this.state.enter, this.state.leave).then(list => {
        this.data = list;
        this.entity.getMap().then(mapData => {
          this.mapData = mapData;
          this.setState({
            loading: false
          });
        });
      });
    });
  }

  onSpaceSelect = (item: Space) => {
    if (item.available) {
      this.setState({
        showConfirm: true,
        selectedSpace: item
      });
    } else if (!item.available && item.bookings && item.bookings.length > 0) {
      this.setState({
        showBookingNames: true,
        selectedSpace: item
      });
    }
  }

  renderItem = (item: Space) => {
    console.log(this.context.showNames);
    const boxStyle: React.CSSProperties = {
      backgroundColor: item.available ? "rgba(48, 209, 88, 0.9)" : "rgba(255, 69, 58, 0.9)",
      position: "absolute",
      left: item.x,
      top: item.y,
      width: item.width,
      height: item.height,
      transform: "rotate: " + item.rotation + "deg",
      cursor: (item.available || (item.bookings && item.bookings.length > 0)) ? "pointer" : "default"
    };
    const textStyle: React.CSSProperties = {
      textAlign: "center"
    };
    const className = (item.width < item.height) ? "space-box space-box-vertical" : "space-box";
    return (
      <div key={item.id} style={boxStyle} className={className} onClick={() => this.onSpaceSelect(item)}>
        <p style={textStyle}>{item.name}</p>
      </div>
    );
  }

  onConfirmBooking = () => {
    this.setState({
      showConfirm: false,
      loading: true
    });
    let booking: Booking = new Booking();
    booking.enter = new Date(this.state.enter);
    booking.leave = new Date(this.state.leave);
    booking.space = this.state.selectedSpace;
    booking.save().then(() => {
      this.setState({
        loading: false,
        showSuccess: true
      });
    });

    this.setState({ showConfirm: false });
  }

  renderBookingNameRow = (booking: Booking) => {
    return (
      <p key={booking.id}>
        {booking.user.email}<br />
        {Formatting.getFormatterShort().format(new Date(booking.enter))}
        &nbsp;&mdash;&nbsp;
        {Formatting.getFormatterShort().format(new Date(booking.leave))}
      </p>
    );
  }

  render() {
    if (this.state.loading) {
      return <Loading />;
    }

    const floorPlanStyle = {
      width: (this.mapData ? this.mapData.width : 0) + "px",
      height: (this.mapData ? this.mapData.height : 0) + "px",
      position: 'relative' as 'relative',
      backgroundImage: (this.mapData ? "url(data:image/" + this.mapData.mapMimeType + ";base64," + this.mapData.data + ")" : "")
    };
    let spaces = this.data.map((item) => {
      return this.renderItem(item);
    });
    return (
      <>
        <div className="container-map">
          <h3>{this.entity.name}</h3>
          <div className="mapScrollContainer">
            <div style={floorPlanStyle}>
              {spaces}
            </div>
          </div>
        </div>
        <Modal show={this.state.showBookingNames} onHide={() => this.setState({ showBookingNames: false })}>
          <Modal.Header closeButton>
            <Modal.Title>{this.state.selectedSpace?.name}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {this.state.selectedSpace?.bookings.map(item => this.renderBookingNameRow(item))}
          </Modal.Body>
          <Modal.Footer>
            <Button variant="primary" onClick={() => this.setState({ showBookingNames: false })}>
              {this.props.t("ok")}
            </Button>
          </Modal.Footer>
        </Modal>
        <Modal show={this.state.showConfirm} onHide={() => this.setState({ showConfirm: false })}>
          <Modal.Header closeButton>
            <Modal.Title>{this.props.t("bookSeat")}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>{this.props.t("space")}: {this.state.selectedSpace?.name}</p>
            <p>{this.props.t("area")}: {this.entity.name}</p>
            <p>{this.props.t("enter")}: {Formatting.getFormatterShort().format(new Date(this.state.enter))}</p>
            <p>{this.props.t("leave")}: {Formatting.getFormatterShort().format(new Date(this.state.leave))}</p>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={() => this.setState({ showConfirm: false })}>
              {this.props.t("cancel")}
            </Button>
            <Button variant="primary" onClick={this.onConfirmBooking}>
              {this.props.t("confirmBooking")}
            </Button>
          </Modal.Footer>
        </Modal>
        <Modal show={this.state.showSuccess} onHide={() => this.setState({ showSuccess: false })} backdrop="static" keyboard={false}>
          <Modal.Header closeButton={false}>
            <Modal.Title>{this.props.t("bookSeat")}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>{this.props.t("bookingConfirmed")}</p>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="primary" as={Link} to="/bookings">
              {this.props.t("myBookings")}
            </Button>
          </Modal.Footer>
        </Modal>
      </>
    )
  }
}

export default withTranslation()(SearchResult as any);
