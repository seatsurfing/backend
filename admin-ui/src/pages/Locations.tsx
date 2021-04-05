import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus } from 'react-feather';
import { Link, Redirect } from 'react-router-dom';
import Loading from '../components/Loading';
import { Location } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
  selectedItem: string
  loading: boolean
}

interface Props {
  t: TFunction
}

class Locations extends React.Component<Props, State> {
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

    let buttons = <Link to="/locations/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> {this.props.t("add")}</Link>;

    if (this.state.loading) {
      return (
        <FullLayout headline={this.props.t("areas")} buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline={this.props.t("areas")} buttons={buttons}>
          <p>{this.props.t("noRecords")}</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline={this.props.t("areas")} buttons={buttons}>
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>{this.props.t("name")}</th>
              <th>{this.props.t("map")}</th>
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

export default withTranslation()(Locations as any);
