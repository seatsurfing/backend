import React from 'react';
import FullLayout from '../components/FullLayout';
import { Form, Col, Row, Button, Alert } from 'react-bootstrap';
import { ChevronLeft as IconBack, Save as IconSave, Trash2 as IconDelete, MapPin as IconMap, Copy as IconCopy } from 'react-feather';
import { Link, RouteChildrenProps, Redirect } from 'react-router-dom';
import Loading from '../components/Loading';
import { Location, Space } from 'flexspace-commons';
import { Rnd } from 'react-rnd';
import './EditLocation.css';

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
  fileLabel: string
  files: FileList | null
  spaces: SpaceState[]
  selectedSpace: number | null
  deleteIds: string[]
  changed: boolean
}

interface Props {
  id: string
}

export default class EditLocation extends React.Component<RouteChildrenProps<Props>, State> {
  entity: Location = new Location();
  mapData: any = null;

  constructor(props: any) {
    super(props);
    this.state = {
      loading: true,
      submitting: false,
      saved: false,
      goBack: false,
      name: "",
      fileLabel: "PNG, JPEG oder GIF",
      files: null,
      spaces: [],
      selectedSpace: null,
      deleteIds: [],
      changed: false
    };
  }

  componentDidMount = () => {
    this.loadData();
  }

  loadData = (id?: string) => {
    if (!id) {
      id = this.props.match?.params.id;
    }
    if (id) {
      Location.get(id).then(location => {
        this.entity = location;
        Space.list(this.entity.id).then(spaces => {
          spaces.forEach(space => {
            this.addRect(space);
          });
          this.entity.getMap().then(mapData => {
            this.mapData = mapData;
            this.setState({
              name: location.name,
              loading: false
            });
          });
        });
      });
    } else {
      this.setState({
        loading: false
      });
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
    this.entity.name = this.state.name;
    this.entity.save().then(() => {
      this.saveSpaces().then(() => {
        if (this.state.files && this.state.files.length > 0) {
          this.entity.setMap(this.state.files.item(0) as File).then(() => {
            this.loadData(this.entity.id);
            this.props.history.push("/locations/" + this.entity.id);
            this.setState({
              saved: true,
              changed: false
            });
          });
        } else {
          this.setState({
            saved: true,
            changed: false
          });
        }
      });
    });
  }

  deleteItem = () => {
    if (window.confirm("Bereich löschen? Alle Plätze und Buchungen gehen verloren!")) {
      this.entity.delete().then(() => {
        this.setState({ goBack: true });
      });
    }
  }

  addRect = (e?: Space): number => {
    let spaces = [...this.state.spaces];
    let space: SpaceState = {
      id: (e ? e.id : ""),
      name: (e ? e.name : "Unbenannt"),
      x: (e ? e.x : 10),
      y: (e ? e.y : 10),
      width: (e ? e.width + "px" : "100px"),
      height: (e ? e.height + "px" : "100px"),
      rotation: 0
    };
    let i = spaces.push(space);
    this.setState({ spaces: spaces, changed: this.state.changed || (e ? false : true) });
    return i;
  }

  setSpacePosition = (i: number, x: number, y: number) => {
    let spaces = [...this.state.spaces];
    let space = { ...spaces[i] };
    space.x = x;
    space.y = y;
    spaces[i] = space;
    this.setState({ spaces: spaces, changed: true });
  }

  setSpaceDimensions = (i: number, width: string, height: string) => {
    let spaces = [...this.state.spaces];
    let space = { ...spaces[i] };
    space.width = width;
    space.height = height;
    spaces[i] = space;
    this.setState({ spaces: spaces, changed: true });
  }

  setSpaceName = (i: number, name: string) => {
    let spaces = [...this.state.spaces];
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
      let spaces = [...this.state.spaces];
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
      let spaces = [...this.state.spaces];
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
      if (!window.confirm("Es gibt ungespeicherte Änderungen! Wirklich verwerfen?")) {
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

  render() {
    if (this.state.goBack) {
      return <Redirect to={`/locations`} />
    }

    let backButton = <Link to="/locations" onClick={this.onBackButtonClick} className="btn btn-sm btn-outline-secondary"><IconBack className="feather" /> Zurück</Link>;
    let buttons = backButton;

    if (this.state.loading) {
      return (
        <FullLayout headline="Bereich bearbeiten" buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let hint = <></>;
    if (this.state.saved) {
      hint = <Alert variant="success">Eintrag wurde aktualisiert.</Alert>
    }

    let buttonDelete = <Button className="btn-sm" variant="outline-secondary" onClick={this.deleteItem}><IconDelete className="feather" /> Löschen</Button>;
    let buttonSave = <Button className="btn-sm" variant="outline-secondary" type="submit" form="form"><IconSave className="feather" /> Speichern</Button>;
    let floorPlan = <></>
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
        buttonCopySpace = <Button className="btn-sm" variant="outline-secondary" onClick={this.copySpace}><IconCopy className="feather" /> Duplizieren</Button>;
        buttonDeleteSpace = <Button className="btn-sm" variant="outline-secondary" onClick={this.deleteSpace}><IconDelete className="feather" /> Platz löschen</Button>;
      }
      floorPlan = (
        <>
          <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 className="h2">Raumplan</h1>
            <div className="btn-toolbar mb-2 mb-md-0">
              <div className="btn-group mr-2">
                {buttonCopySpace} {buttonDeleteSpace}
                <Button className="btn-sm" variant="outline-secondary" onClick={() => this.addRect()}><IconMap className="feather" /> Platz hinzufügen</Button>
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
    } else {
      buttons = <>{backButton} {buttonSave}</>;
    }
    return (
      <FullLayout headline="Bereich bearbeiten" buttons={buttons}>
        <Form onSubmit={this.onSubmit} id="form">
          {hint}
          <Form.Group as={Row}>
            <Form.Label column sm="2">Name</Form.Label>
            <Col sm="4">
              <Form.Control type="text" placeholder="Name" value={this.state.name} onChange={(e: any) => this.setState({ name: e.target.value })} required={true} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Raumplan</Form.Label>
            <Col sm="4">
              <Form.File label={this.state.fileLabel} custom={true} accept="image/png, image/jpeg, image/gif" onChange={(e: any) => this.setState({ files: e.target.files, fileLabel: e.target.files.item(0).name })} required={!this.entity.id} />
            </Col>
          </Form.Group>
        </Form>
        {floorPlan}
      </FullLayout>
    );
  }
}
