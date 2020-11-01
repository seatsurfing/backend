import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus } from 'react-feather';
import { Link, Redirect } from 'react-router-dom';
import Loading from '../components/Loading';
import { Location } from 'flexspace-commons';

interface State {
  selectedItem: string
  loading: boolean
}

export default class Locations extends React.Component<{}, State> {
  data: Location[] = [];

  constructor(props: any) {
    super(props);
    this.state = {
      selectedItem: "",
      loading: true
    };
  }
  
  componentDidMount = () => {
    this.loadItems();
  }

  loadItems = () => {
    Location.list().then(list => {
      this.data = list;
      this.setState({ loading: false });
    });
  }

  onItemSelect = (location: Location) => {
    this.setState({ selectedItem: location.id });
  }

  renderItem = (location: Location) => {
    return (
      <tr key={location.id} onClick={() => this.onItemSelect(location)}>
        <td>{location.name}</td>
        <td>{location.mapWidth}x{location.mapHeight}</td>
      </tr>
    );
  }

  render() {
    if (this.state.selectedItem) {
      return <Redirect to={`/locations/${this.state.selectedItem}`} />
    }

    let buttons = <Link to="/locations/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> Neu</Link>;

    if (this.state.loading) {
      return (
        <FullLayout headline="Bereiche" buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline="Bereiche" buttons={buttons}>
          <p>Keine Datensätze gefunden.</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline="Bereiche" buttons={buttons}>
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Karte</th>
            </tr>
          </thead>
          <tbody>
            {rows}
          </tbody>
        </Table>
      </FullLayout>
    );
  }
}
