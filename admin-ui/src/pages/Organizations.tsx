import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus } from 'react-feather';
import { Link, Navigate } from 'react-router-dom';
import Loading from '../components/Loading';
import { Organization } from 'flexspace-commons';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
  selectedItem: string
  loading: boolean
}

interface Props {
  t: TFunction
}

class Organizations extends React.Component<Props, State> {
  data: Organization[] = [];

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
    Organization.list().then(list => {
      this.data = list;
      this.setState({ loading: false });
    });
  }

  onItemSelect = (org: Organization) => {
    this.setState({ selectedItem: org.id });
  }

  renderItem = (org: Organization) => {
    return (
      <tr key={org.id} onClick={() => this.onItemSelect(org)}>
        <td>{org.name}</td>
      </tr>
    );
  }

  render() {
    if (this.state.selectedItem) {
      return <Navigate replace={true} to={`/organizations/${this.state.selectedItem}`} />
    }

    let buttons = <Link to="/organizations/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> {this.props.t("add")}</Link>;

    if (this.state.loading) {
      return (
        <FullLayout headline={this.props.t("organizations")} buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline={this.props.t("organizations")} buttons={buttons}>
          <p>{this.props.t("noRecords")}</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline={this.props.t("organizations")} buttons={buttons}>
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>{this.props.t("org")}</th>
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

export default withTranslation()(Organizations as any);
