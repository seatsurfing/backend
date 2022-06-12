import React, { RefObject } from 'react';
import './Login.css';
import { Form, Alert, Col, Row, Modal, Button, ListGroup } from 'react-bootstrap';
import { Location, Booking, Ajax, Formatting, Space, AjaxError, UserPreference } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';
// @ts-ignore
import DateTimePicker from 'react-datetime-picker';
import DatePicker from 'react-date-picker';
import './Search.css';
import { AuthContext } from '../AuthContextData';
import Loading from '../components/Loading';
import { EnterOutline as EnterIcon, ExitOutline as ExitIcon, LocationOutline as LocationIcon, ChevronUpOutline as CollapseIcon, SettingsOutline as SettingsIcon, MapOutline as MapIcon, ListOutline as ListIcon } from 'react-ionicons'
import ErrorText from '../types/ErrorText';
import { NavigateFunction } from 'react-router-dom';
import { withNavigate } from '../types/withNavigate';

interface State {
  enter: Date
  leave: Date
  locationId: string
  canSearch: boolean
  canSearchHint: string
  showBookingNames: boolean
  selectedSpace: Space | null
  showConfirm: boolean
  showSuccess: boolean
  showError: boolean
  errorText: string
  loading: boolean
  listView: boolean
  prefEnterTime: number
  prefWorkdayStart: number
  prefWorkdayEnd: number
  prefWorkdays: number[]
  prefLocationId: string
}

interface Props {
  navigate: NavigateFunction
  t: TFunction
}

class Search extends React.Component<Props, State> {
  static contextType = AuthContext;
  static PreferenceEnterTimeNow: number = 1;
	static PreferenceEnterTimeNextDay: number = 2;
	static PreferenceEnterTimeNextWorkday: number = 3;

  data: Space[];
  locations: Location[]
  mapData: any;
  curBookingCount: number = 0;
  searchContainerRef: RefObject<any>;
  enterChangeTimer: number | undefined;
  leaveChangeTimer: number | undefined;

  constructor(props: any) {
    super(props);
    this.data = [];
    this.locations = [];
    this.mapData = null;
    this.searchContainerRef = React.createRef();
    this.enterChangeTimer = undefined;
    this.leaveChangeTimer = undefined;
    this.state = {
      enter: new Date(),
      leave: new Date(),
      locationId: "",
      canSearch: false,
      canSearchHint: "",
      showBookingNames: false,
      selectedSpace: null,
      showConfirm: false,
      showSuccess: false,
      showError: false,
      errorText: "",
      loading: true,
      listView: false,
      prefEnterTime: 0,
      prefWorkdayStart: 0,
      prefWorkdayEnd: 0,
      prefWorkdays: [],
      prefLocationId: "",
    };
  }

  componentDidMount = () => {
    this.loadItems();
  }

  loadItems = () => {
    let promises = [
      this.loadLocations(),
      this.loadPreferences(),
    ];
    Promise.all(promises).then(() => {
      this.initDates();
      if (this.state.locationId === "" && this.locations.length > 0) {
        let defaultLocationId = this.locations[0].id;
        if (this.state.prefLocationId) {
          this.locations.forEach(location => {
            if (location.id === this.state.prefLocationId) {
              defaultLocationId = this.state.prefLocationId;
            }
          })
        }
        this.setState({ locationId: defaultLocationId });
        this.loadMap(this.state.locationId).then(() => {
          this.setState({ loading: false });
        });
      } else {
        this.setState({ loading: false });
      }
    });
  }

  loadPreferences = async (): Promise<void> => {
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      UserPreference.list().then(list => {
        let state: any = {};
        list.forEach(s => {
          if (s.name === "enter_time") state.prefEnterTime = window.parseInt(s.value);
          if (s.name === "workday_start") state.prefWorkdayStart = window.parseInt(s.value);
          if (s.name === "workday_end") state.prefWorkdayEnd = window.parseInt(s.value);
          if (s.name === "workdays") state.prefWorkdays = s.value.split(",").map(val => window.parseInt(val));
          if (s.name === "location_id") state.prefLocationId = s.value;
        });
        if (self.context.dailyBasisBooking) {
          state.prefWorkdayStart = 0;
          state.prefWorkdayEnd = 23;
        }
        self.setState({
          ...state
        }, () => resolve());
      }).catch(e => reject(e));
    });
  }

  initCurrentBookingCount = () => {
    Booking.list().then(list => {
      this.curBookingCount = list.length;
      this.updateCanSearch();
    });
  }

  initDates = () => {
    let enter = new Date();
    if (this.state.prefEnterTime === Search.PreferenceEnterTimeNow) {
      enter.setHours(enter.getHours() + 1, 0, 0);
      if (enter.getHours() < this.state.prefWorkdayStart) {
        enter.setHours(this.state.prefWorkdayStart, 0, 0, 0);
      }
      if (enter.getHours() >= this.state.prefWorkdayEnd) {
        enter.setDate(enter.getDate() + 1);
        enter.setHours(this.state.prefWorkdayStart, 0, 0, 0);
      }
    } else if (this.state.prefEnterTime === Search.PreferenceEnterTimeNextDay) {
      enter.setDate(enter.getDate() + 1);
      enter.setHours(this.state.prefWorkdayStart, 0, 0, 0);
    } else if (this.state.prefEnterTime === Search.PreferenceEnterTimeNextWorkday) {
      enter.setDate(enter.getDate() + 1);
      let add = 0;
      let nextDayFound = false;
      let lookFor = enter.getDay();
      while (!nextDayFound) {
        if (this.state.prefWorkdays.includes(lookFor) || add > 7) {
          nextDayFound = true;
        } else {
          add++;
          lookFor++;
          if (lookFor > 6) {
            lookFor = 0;
          }
        }
      }
      enter.setDate(enter.getDate() + add);
      enter.setHours(this.state.prefWorkdayStart, 0, 0, 0);
    }

    let leave = new Date(enter);
    leave.setHours(this.state.prefWorkdayEnd, 0, 0);

    if (this.context.dailyBasisBooking) {
      enter.setHours(0, 0, 0, 0);
      leave.setHours(23, 59, 59, 0);
    }

    this.setState({
      enter: enter,
      leave: leave
    });
  }

  loadLocations = async (): Promise<void> => {
    return Location.list().then(list => {
      this.locations = list;
    });
  }

  loadMap = async (locationId: string) => {
    this.setState({ loading: true });
    return Location.get(locationId).then(location => {
      return this.loadSpaces(location.id).then(() => {
        return Ajax.get(location.getMapUrl()).then(mapData => {
          this.mapData = mapData.json;
          this.centerMapView();
        });
      });
    })
  }

  centerMapView = () => {
    let timer: number | undefined = undefined;
    let cb = () => {
      const el = document.querySelector('.mapScrollContainer');
      if (el) {
        window.clearInterval(timer);
        el.scrollLeft = (this.mapData ? this.mapData.width : 0) / 2 - (window.innerWidth / 2);
        el.scrollTop = (this.mapData ? this.mapData.height : 0) / 2 - (window.innerHeight / 2);
      }
    };
    timer = window.setInterval(cb, 10);
  }

  loadSpaces = async (locationId: string) => {
    this.setState({ loading: true });
    return Space.listAvailability(locationId, this.state.enter, this.state.leave).then(list => {
      this.data = list;
    });
  }

  updateCanSearch = async () => {
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
    let self = this;
    return new Promise<void>(function (resolve, reject) {
      self.setState({
        canSearch: res,
        canSearchHint: hint
      }, () => resolve());
    });
  }

  renderLocations = () => {
    return this.locations.map(location => {
      return <option value={location.id} key={location.id}>{location.name}</option>;
    });
  }

  setEnterDate = (value: Date | Date[]) => {
    let dateChangedCb = () => {
      this.updateCanSearch().then(() => {
        if (!this.state.canSearch) {
          this.setState({ loading: false });
        } else {
          let promises = [
            this.initCurrentBookingCount(),
            this.loadSpaces(this.state.locationId),
          ];
          Promise.all(promises).then(() => {
            this.setState({ loading: false });
          });
        }
      });
    };
    let performChange = () => {
      let diff = this.state.leave.getTime() - this.state.enter.getTime();
      let date = (value instanceof Date) ? value : value[0];
      if (this.context.dailyBasisBooking) {
        date.setHours(0, 0, 0);
      }
      let leave = new Date();
      leave.setTime(date.getTime() + diff);
      this.setState({
        enter: date,
        leave: leave
      }, () => dateChangedCb());
    };
    window.clearTimeout(this.enterChangeTimer);
    this.enterChangeTimer = window.setTimeout(performChange, 2000);
  }

  setLeaveDate = (value: Date | Date[]) => {
    let dateChangedCb = () => {
      this.updateCanSearch().then(() => {
        if (!this.state.canSearch) {
          this.setState({ loading: false });
        } else {
          let promises = [
            this.initCurrentBookingCount(),
            this.loadSpaces(this.state.locationId),
          ];
          Promise.all(promises).then(() => {
            this.setState({ loading: false });
          });
        }
      });
    };
    let performChange = () => {
      let date = (value instanceof Date) ? value : value[0];
      if (this.context.dailyBasisBooking) {
        date.setHours(23, 59, 59);
      }
      this.setState({
        leave: date
      }, () => dateChangedCb());
    };
    window.clearTimeout(this.leaveChangeTimer);
    this.leaveChangeTimer = window.setTimeout(performChange, 2000);
  }

  changeLocation = (id: string) => {
    this.setState({
      locationId: id,
      loading: true,
    });
    this.loadMap(id).then(() => {
      this.setState({ loading: false });
    });
  }

  onSpaceSelect = (item: Space) => {
    if (item.available) {
      this.setState({
        showConfirm: true,
        selectedSpace: item
      });
    } else {
      let bookings = Booking.createFromRawArray(item.rawBookings);
      if (!item.available && bookings && bookings.length > 0) {
        this.setState({
          showBookingNames: true,
          selectedSpace: item
        });
      }
    }
  }

  getAvailibilityStyle = (item: Space, bookings: Booking[]) => {
    const mydesk = (bookings.find(b => b.user.email === this.context.username ));
    return (mydesk ? " space-mydesk" : (item.available ? " space-available" : " space-notavailable"));
  }
  
  getBookersList = (bookings: Booking[]) => {
    if (!bookings.length) return "";
    let str = "";
    bookings.forEach(b => {
      str += (str ? ", " : "") + b.user.email
     });
     return str;
  }

  renderItem = (item: Space) => {
    let bookings = Booking.createFromRawArray(item.rawBookings);
    const boxStyle: React.CSSProperties = {
      position: "absolute",
      left: item.x,
      top: item.y,
      width: item.width,
      height: item.height,
      transform: "rotate: " + item.rotation + "deg",
      cursor: (item.available || (bookings && bookings.length > 0)) ? "pointer" : "default"
    };
    const textStyle: React.CSSProperties = {
      textAlign: "center"
    };
    const className = "space space-box" 
      + ((item.width < item.height) ? " space-box-vertical" : "")
      + this.getAvailibilityStyle(item, bookings);
    return (
      <div key={item.id} style={boxStyle} className={className}
        onClick={() => this.onSpaceSelect(item)}
        title={this.getBookersList(bookings)}>
        <p style={textStyle}>{item.name}</p>
      </div>
    );
  }

  renderListItem = (item: Space) => {
    let bookings: Booking[] = [];
    bookings = Booking.createFromRawArray(item.rawBookings);
    const className = "space space-listitem" + this.getAvailibilityStyle(item, bookings);
    return (
      <ListGroup.Item key={item.id} action={true} onClick={(e) => { e.preventDefault(); this.onSpaceSelect(item); }}>
        <p className={className}>
          <span key={item.id} className="space-listitem-title">{item.name}</span>
          {bookings.map((booking, index) => <span key={booking.user.id}><br/>{booking.user.email}</span>)}
        </p>
      </ListGroup.Item>
    );
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

  onConfirmBooking = () => {
    if (this.state.selectedSpace == null) {
      return;
    }
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
    }).catch(e => {
      let code: number = 0;
      if (e instanceof AjaxError) {
        code = e.appErrorCode;
      }
      this.setState({
        loading: false,
        showError: true,
        errorText: ErrorText.getTextForAppCode(code, this.props.t, this.context)
      });
    });
  }

  getLocationName = (): string => {
    let name: string = this.props.t("none");
    this.locations.forEach(location => {
      if (this.state.locationId === location.id) {
        name = location.name;
      }
    });
    return name;
  }

  toggleSearchContainer = () => {
    const ref = this.searchContainerRef.current;
    ref.classList.toggle("minimized");
  }

  toggleListView = () => {
    this.setState({ listView: !this.state.listView });
  }

  render() {
    let hint = <></>;
    if ((!this.state.canSearch) && (this.state.canSearchHint)) {
      hint = (
        <Form.Group>
          <Alert variant="warning">{this.state.canSearchHint}</Alert>
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

    let listOrMap = <></>;
    if (this.state.listView) {
      listOrMap = (
        <div className="container-signin">
          <Form className="form-signin">
            <ListGroup>
              {this.data.map(item => this.renderListItem(item))}
            </ListGroup>
          </Form>
        </div>
      );
    } else {
      const floorPlanStyle = {
        width: (this.mapData ? this.mapData.width : 0) + "px",
        height: (this.mapData ? this.mapData.height : 0) + "px",
        position: 'relative' as 'relative',
        backgroundImage: (this.mapData ? "url(data:image/" + this.mapData.mapMimeType + ";base64," + this.mapData.data + ")" : "")
      };
      let spaces = this.data.map((item) => {
        return this.renderItem(item);
      });
      listOrMap = (
        <div className="container-map">
          <div className="mapScrollContainer">
            <div style={floorPlanStyle}>
              {spaces}
            </div>
          </div>
        </div>
      );
    }

    let configContainer = (
      <div className="container-search-config" ref={this.searchContainerRef}>
        <div className="collapse-bar" onClick={() => this.toggleSearchContainer()}>
          <CollapseIcon color={'#000'} height="20px" width="20px" cssClasses="collapse-icon" />
          <SettingsIcon color={'#555'} height="26px" width="26px" cssClasses="expand-icon" />
        </div>
        <div className="content">
          <Form>
            <Form.Group as={Row}>
              <Col xs="2"><LocationIcon title={this.props.t("area")} color={'#555'} height="20px" width="20px" /></Col>
              <Col xs="10">
                <Form.Select required={true} value={this.state.locationId} onChange={(e) => this.changeLocation(e.target.value)}>
                  {this.renderLocations()}
                </Form.Select>
              </Col>
            </Form.Group>
            <Form.Group as={Row} className="margin-top-10">
              <Col xs="2"><EnterIcon title={this.props.t("enter")} color={'#555'} height="20px" width="20px" /></Col>
              <Col xs="10">
                {enterDatePicker}
              </Col>
            </Form.Group>
            <Form.Group as={Row} className="margin-top-10">
              <Col xs="2"><ExitIcon title={this.props.t("leave")} color={'#555'} height="20px" width="20px" /></Col>
              <Col xs="10">
                {leaveDatePicker}
              </Col>
            </Form.Group>
            {hint}
          </Form>
          <Button variant="outline-dark" onClick={() => this.toggleListView()} className="margin-top-10">
            <MapIcon title={this.props.t("showMap")} color={'#555'} height="26px" width="26px" style={{ "display": this.state.listView ? "" : "none" }} />
            <ListIcon title={this.props.t("showList")} color={'#555'} height="26px" width="26px" style={{ "display": this.state.listView ? "none" : "" }} />
          </Button>
        </div>
      </div>
    );

    let confirmModal = (
      <Modal show={this.state.showConfirm} onHide={() => this.setState({ showConfirm: false })}>
        <Modal.Header closeButton>
          <Modal.Title>{this.props.t("bookSeat")}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>{this.props.t("space")}: {this.state.selectedSpace?.name}</p>
          <p>{this.props.t("area")}: {this.getLocationName()}</p>
          <p>{this.props.t("enter")}: {Formatting.getFormatterShort().format(Formatting.convertToFakeUTCDate(new Date(this.state.enter)))}</p>
          <p>{this.props.t("leave")}: {Formatting.getFormatterShort().format(Formatting.convertToFakeUTCDate(new Date(this.state.leave)))}</p>
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
    );
    let bookings: Booking[] = [];
    if (this.state.selectedSpace) {
      bookings = Booking.createFromRawArray(this.state.selectedSpace.rawBookings);
    }
    const myBooking = (bookings.find(b => b.user.email === this.context.username ));
    let gotoBooking;
    if (myBooking){
      gotoBooking = (
        <Button variant="secondary" onClick={() => {
          this.setState({ showBookingNames: false })
          this.props.navigate("/bookings#"+myBooking.id)
        }}>
          {this.props.t("gotoBooking")}
        </Button>
      )
    }
    let bookingNamesModal = (
      <Modal show={this.state.showBookingNames} onHide={() => this.setState({ showBookingNames: false })}>
        <Modal.Header closeButton>
          <Modal.Title>{this.state.selectedSpace?.name}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {bookings.map(item => {
            return <span key={item.user.id}>{this.renderBookingNameRow(item)}</span>
          })}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary" onClick={() => this.setState({ showBookingNames: false })}>
            {this.props.t("ok")}
          </Button>
        {gotoBooking}
        </Modal.Footer>
      </Modal>
    );
    let successModal = (
      <Modal show={this.state.showSuccess} onHide={() => this.setState({ showSuccess: false })} backdrop="static" keyboard={false}>
        <Modal.Header closeButton={false}>
          <Modal.Title>{this.props.t("bookSeat")}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>{this.props.t("bookingConfirmed")}</p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary" onClick={() => this.props.navigate("/bookings")}>
            {this.props.t("myBookings").toString()}
          </Button>
          <Button variant="secondary" onClick={() => {
            this.setState({ showSuccess: false });
            this.refreshPage();
            }}>
            {this.props.t("ok").toString()}
          </Button>
        </Modal.Footer>
      </Modal>
    );
    let errorModal = (
      <Modal show={this.state.showError} onHide={() => this.setState({ showError: false })} backdrop="static" keyboard={false}>
        <Modal.Header closeButton={false}>
          <Modal.Title>{this.props.t("bookSeat")}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>{this.state.errorText}</p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary" onClick={() => this.props.navigate("/bookings")}>
            {this.props.t("myBookings").toString()}
          </Button>
        </Modal.Footer>
      </Modal>
    );

    if (this.state.loading) {
      return (
        <>
          <Loading />
          {configContainer}
        </>
      );
    }

    return (
      <>
        {confirmModal}
        {bookingNamesModal}
        {successModal}
        {errorModal}
        {listOrMap}
        {configContainer}
      </>
    )
  }

  refreshPage = () => {
    this.setState({
      loading: true,
    });
    this.loadMap(this.state.locationId).then(() => {
      this.setState({ loading: false });
    });
  }
}

export default withNavigate(withTranslation()(Search as any));
