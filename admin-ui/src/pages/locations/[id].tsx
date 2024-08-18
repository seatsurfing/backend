import React, { useEffect } from 'react';
import FullLayout from '../../components/FullLayout';
import { Form, Col, Row, Button, Alert, InputGroup, Table } from 'react-bootstrap';
import { ChevronLeft as IconBack, Save as IconSave, Trash2 as IconDelete, MapPin as IconMap, Copy as IconCopy, Loader as IconLoad, Download as IconDownload } from 'react-feather';
import Loading from '../../components/Loading';
import { Ajax, Location, Space } from 'flexspace-commons';
import { Rnd } from 'react-rnd';
import { WithTranslation, withTranslation } from 'next-i18next';
import { NextRouter } from 'next/router';
import Link from 'next/link';
import withReadyRouter from '@/components/withReadyRouter';

interface SpaceState {
  id: string
  name: string
  x: number
  y: number
  width: string
  height: string
  rotation: number
}

interface State {
  loading: boolean
  submitting: boolean
  saved: boolean
  goBack: boolean
  name: string
  description: string
  limitConcurrentBookings: boolean
  maxConcurrentBookings: number
  timezone: string
  fileLabel: string
  files: FileList | null
  spaces: SpaceState[]
  selectedSpace: number | null
  deleteIds: string[]
  changed: boolean
}

interface Props extends WithTranslation {
  router: NextRouter
}

class EditLocation extends React.Component<Props, State> {
  entity: Location = new Location();
  mapData: any = null;
  timezones: string[];
  ExcellentExport: any;

  constructor(props: any) {
    super(props);
    this.timezones = [];
    this.state = {
      loading: true,
      submitting: false,
      saved: false,
      goBack: false,
      name: "",
      description: "",
      limitConcurrentBookings: false,
      maxConcurrentBookings: 0,
      timezone: "",
      fileLabel: this.props.t("mapFileTypes"),
      files: null,
      spaces: [],
      selectedSpace: null,
      deleteIds: [],
      changed: false
    };
  }

  componentDidMount = () => {
    if (!Ajax.CREDENTIALS.accessToken) {
      this.props.router.push("/login");
      return;
    }
    let promises = [
      this.loadData(),
      this.loadTimezones(),
    ];
    Promise.all(promises).then(() => {
      this.setState({
        loading: false
      });
    });
    import('excellentexport').then(imp => this.ExcellentExport = imp.default);
  }

  loadTimezones = async (): Promise<void> => {
    return Ajax.get("/setting/timezones").then(res => {
      this.timezones = res.json;
    });
  }

  loadData = async (locationId?: string): Promise<void> => {
    if (!locationId) {
      const { id } = this.props.router.query;
      if (id && (typeof id === "string") && (id !== 'add')) {
        locationId = id;
      }
    }
    if (locationId) {
      return Location.get(locationId).then(location => {
        this.entity = location;
        return Space.list(this.entity.id).then(spaces => {
          this.setState({ spaces: spaces.map( (s) => this.newSpaceState(s) ) });
          return this.entity.getMap().then(mapData => {
            this.mapData = mapData;
            this.setState({
              name: location.name,
              description: location.description,
              limitConcurrentBookings: (location.maxConcurrentBookings > 0),
              maxConcurrentBookings: location.maxConcurrentBookings,
              timezone: location.timezone,
              loading: false
            });
          });
        });
      });
    } else {
      //return Promise.resolve();
    }
  }

  saveSpaces = async () => {
    let spaces = await Space.list(this.entity.id);
    for (let space of spaces) {
      if (this.state.deleteIds.indexOf(space.id) > -1) {
        await space.delete();
      }
    }
    for (let item of this.state.spaces) {
      let space: Space = new Space();
      spaces.forEach((spaceItem) => {
        if (item.id === spaceItem.id) {
          space = spaceItem
        }
      });
      space.locationId = this.entity.id;
      space.name = item.name;
      space.x = item.x;
      space.y = item.y;
      space.width = parseInt(item.width.replace(/^\D+/g, ''));
      space.height = parseInt(item.height.replace(/^\D+/g, ''));
      space.rotation = item.rotation;
      await space.save();
    }
  }

  onSubmit = (e: any) => {
    e.preventDefault();
    this.setState({ submitting: true });
    this.entity.name = this.state.name;
    this.entity.description = this.state.description;
    this.entity.maxConcurrentBookings = (this.state.limitConcurrentBookings ? this.state.maxConcurrentBookings : 0);
    this.entity.timezone = this.state.timezone;
    this.entity.save().then(() => {
      this.saveSpaces().then(() => {
        if (this.state.files && this.state.files.length > 0) {
          this.entity.setMap(this.state.files.item(0) as File).then(() => {
            this.loadData(this.entity.id);
            this.props.router.push("/locations/" + this.entity.id);
            this.setState({
              files: null,
              saved: true,
              changed: false,
              submitting: false
            });
          });
        } else {
          this.setState({
            saved: true,
            changed: false,
            submitting: false
          });
        }
      });
    });
  }

  deleteItem = () => {
    if (window.confirm(this.props.t("confirmDeleteArea"))) {
      this.entity.delete().then(() => {
        this.setState({ goBack: true });
      });
    }
  }

  newSpaceState = (e?: Space): SpaceState => {
    return {
      id: (e ? e.id : ""),
      name: (e ? e.name : this.props.t("unnamed")),
      x: (e ? e.x : 10),
      y: (e ? e.y : 10),
      width: (e ? e.width + "px" : "100px"),
      height: (e ? e.height + "px" : "100px"),
      rotation: 0
    };
  }

  addRect = (e?: Space): number => {
    let spaces = this.state.spaces;
    let space = this.newSpaceState(e);
    let i = spaces.push(space);
    this.setState({ spaces: spaces, changed: this.state.changed || (e ? false : true) });
    return i;
  }

  setSpacePosition = (i: number, x: number, y: number) => {
    let spaces = this.state.spaces;
    let space = { ...spaces[i] };
    space.x = x;
    space.y = y;
    spaces[i] = space;
    this.setState({ spaces: spaces, changed: true });
  }

  setSpaceDimensions = (i: number, width: string, height: string) => {
    let spaces = this.state.spaces;
    let space = { ...spaces[i] };
    space.width = width;
    space.height = height;
    spaces[i] = space;
    this.setState({ spaces: spaces, changed: true });
  }

  setSpaceName = (i: number, name: string) => {
    let spaces = this.state.spaces;
    let space = { ...spaces[i] };
    space.name = name;
    spaces[i] = space;
    this.setState({ spaces: spaces, changed: true });
  }

  onSpaceSelect = (i: number) => {
    this.setState({ selectedSpace: i });
  }

  copySpace = () => {
    if (this.state.selectedSpace != null) {
      let spaces = this.state.spaces;
      let space = { ...spaces[this.state.selectedSpace] };
      let newSpace: SpaceState = Object.assign({}, space);
      newSpace.id = "";
      newSpace.x += 20;
      newSpace.y += 20;
      spaces.push(newSpace);
      this.setState({ spaces: spaces });
      this.setState({ selectedSpace: null, changed: true });
    }
  }

  deleteSpace = () => {
    if (this.state.selectedSpace != null) {
      let spaces = this.state.spaces;
      let space = { ...spaces[this.state.selectedSpace] };
      if (space.id) {
        let deleteIds = [...this.state.deleteIds];
        deleteIds.push(space.id);
        this.setState({ deleteIds: deleteIds });
      }
      spaces.splice(this.state.selectedSpace, 1);
      this.setState({ spaces: spaces });
      this.setState({ selectedSpace: null, changed: true });
    }
  }

  onBackButtonClick = (e: any) => {
    if (this.state.changed) {
      if (!window.confirm(this.props.t("confirmDiscard"))) {
        e.preventDefault();
      }
    }
  }

  renderRect = (i: number) => {
    let size = { width: this.state.spaces[i].width, height: this.state.spaces[i].height };
    let position = { x: this.state.spaces[i].x, y: this.state.spaces[i].y };
    let width = parseInt(this.state.spaces[i].width.replace(/^\D+/g, ''));
    let height = parseInt(this.state.spaces[i].height.replace(/^\D+/g, ''));
    let className = "space-dragger";
    let inputStyle = {};
    if (width < height) {
      className += " space-dragger-vertical";
      inputStyle = {
        width: height + "px"
      };
    }
    if (i === this.state.selectedSpace) {
      className += " space-dragger-selected";
    }
    return <Rnd
      key={i}
      size={size}
      position={position}
      onDragStop={(e, d) => { this.setSpacePosition(i, d.x, d.y); this.onSpaceSelect(i); }}
      onResizeStop={(e, d, ref) => { this.setSpaceDimensions(i, ref.style.width, ref.style.height) }}
      className={className}>
      <input type="text" style={inputStyle} value={this.state.spaces[i].name} onChange={(e) => { this.setSpaceName(i, e.target.value) }} />
    </Rnd>;
  }

  getSaveButton = () => {
    if (this.state.submitting) {
      return <Button className="btn-sm" variant="outline-secondary" type="submit" form="form" disabled={true}><IconLoad className="feather loader" /> {this.props.t("save")}</Button>;
    } else {
      return <Button className="btn-sm" variant="outline-secondary" type="submit" form="form"><IconSave className="feather" /> {this.props.t("save")}</Button>;
    }
  }

  renderRow = (space: SpaceState) => {
    return (
      <tr key={space.id} >
        <td>{space.name}</td>
        <td>{window.location.origin}/ui/search?lid={this.entity.id}&sid={space.id}</td>
      </tr>
    );
  }

  exportTable = (e: any) => {
    return this.ExcellentExport.convert(
      { anchor: e.target, filename: "seatsurfing-spaces", format: "xlsx" },
      [{ name: "Seatsurfing Spaces", from: { table: "datatable" } }]
    );
  }

  render() {
    if (this.state.goBack) {
      this.props.router.push(`/locations`);
      return <></>
    }

    let backButton = <Link href="/locations" onClick={this.onBackButtonClick} className="btn btn-sm btn-outline-secondary"><IconBack className="feather" /> {this.props.t("back")}</Link>;
    let buttons = backButton;

    if (this.state.loading) {
      return (
        <FullLayout headline={this.props.t("editArea")} buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let hint = <></>;
    if (this.state.saved) {
      hint = <Alert variant="success">{this.props.t("entryUpdated")}</Alert>
    }

    let buttonDelete = <Button className="btn-sm" variant="outline-secondary" onClick={this.deleteItem}><IconDelete className="feather" /> {this.props.t("delete")}</Button>;
    let buttonSave = this.getSaveButton();
    let floorPlan = <></>
    let spaceTable = <></>
    let rows = this.state.spaces.map((item) => this.renderRow(item));
    if (this.entity.id) {
      buttons = <>{backButton} {buttonDelete} {buttonSave}</>;
      const floorPlanStyle = {
        width: (this.mapData ? this.mapData.width : 0) + "px",
        height: (this.mapData ? this.mapData.height : 0) + "px",
        position: 'relative' as 'relative',
        backgroundImage: (this.mapData ? "url(data:image/" + this.mapData.mapMimeType + ";base64," + this.mapData.data + ")" : "")
      };
      let spaces = this.state.spaces.map((item, i) => {
        return this.renderRect(i);
      });
      let buttonCopySpace = <></>;
      let buttonDeleteSpace = <></>;
      if (this.state.selectedSpace != null) {
        buttonCopySpace = <Button className="btn-sm" variant="outline-secondary" onClick={this.copySpace}><IconCopy className="feather" /> {this.props.t("duplicate")}</Button>;
        buttonDeleteSpace = <Button className="btn-sm" variant="outline-secondary" onClick={this.deleteSpace}><IconDelete className="feather" /> {this.props.t("deleteSpace")}</Button>;
      }
      floorPlan = (
        <>
          <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 className="h2">{this.props.t("floorplan")}</h1>
            <div className="btn-toolbar mb-2 mb-md-0">
              <div className="btn-group me-2">
                {buttonCopySpace} {buttonDeleteSpace}
                <Button className="btn-sm" variant="outline-secondary" onClick={() => this.addRect()}><IconMap className="feather" /> {this.props.t("addSpace")}</Button>
              </div>
            </div>
          </div>
          <div className="mapScrollContainer">
            <div style={floorPlanStyle}>
              {spaces}
            </div>
          </div>
        </>
      );
      let downloadButton = <a download={`seatsurfing-${this.state.name}-spaces.xlsx`} href="#" className="btn btn-sm btn-outline-secondary" onClick={this.exportTable}><IconDownload className="feather" /> {this.props.t("download")}</a>;
      spaceTable = (
        <>
          <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">{this.props.t("spaces")}</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <div className="btn-group me-2">
              {downloadButton}
            </div>
          </div>
        </div>
        <Table striped={true} hover={true} id="datatable">
          <thead>
            <tr>
              <th>{this.props.t("name")}</th>
              <th>{this.props.t("bookingLink")}</th>
            </tr>
          </thead>
          <tbody>
            {rows}
          </tbody>
        </Table>
      </>
      );
    } else {
      buttons = <>{backButton} {buttonSave}</>;
    }
    return (
      <FullLayout headline={this.props.t("editArea")} buttons={buttons}>
        <Form onSubmit={this.onSubmit} id="form">
          {hint}
          <Form.Group as={Row}>
            <Form.Label column sm="2">{this.props.t("name")}</Form.Label>
            <Col sm="4">
              <Form.Control type="text" placeholder={this.props.t("name")} value={this.state.name} onChange={(e: any) => this.setState({ name: e.target.value })} required={true} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">{this.props.t("description")}</Form.Label>
            <Col sm="4">
              <Form.Control type="text" placeholder={this.props.t("description")} value={this.state.description} onChange={(e: any) => this.setState({ description: e.target.value })} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">{this.props.t("timezone")}</Form.Label>
            <Col sm="4">
              <Form.Select value={this.state.timezone} onChange={(e: any) => this.setState({ timezone: e.target.value })}>
                <option value="">({this.props.t("default")})</option>
                {this.timezones.map(tz => <option key={tz} value={tz}>{tz}</option>)}
              </Form.Select>
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">{this.props.t("maxConcurrentBookings")}</Form.Label>
            <Col sm="4">
              <InputGroup>
                <InputGroup.Checkbox type="checkbox" id="check-limitConcurrentBookings" checked={this.state.limitConcurrentBookings} onChange={(e: any) => this.setState({ limitConcurrentBookings: e.target.checked })} />
                <Form.Control type="number" min="0" value={this.state.maxConcurrentBookings} onChange={(e: any) => this.setState({ maxConcurrentBookings: parseInt(e.target.value) })} disabled={!this.state.limitConcurrentBookings} />
              </InputGroup>
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">{this.props.t("floorplan")}</Form.Label>
            <Col sm="4">
              <Form.Control type="file" accept="image/png, image/jpeg, image/gif" onChange={(e: any) => this.setState({ files: e.target.files, fileLabel: e.target.files.item(0).name })} required={!this.entity.id} />
            </Col>
          </Form.Group>
        </Form>
        {floorPlan}
        {spaceTable}
      </FullLayout>
    );
  }
}

export default withTranslation(['admin'])(withReadyRouter(EditLocation as any));
