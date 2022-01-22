import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus, Download as IconDownload } from 'react-feather';
import { Link, Navigate } from 'react-router-dom';
import Loading from '../components/Loading';
import { User, AuthProvider } from 'flexspace-commons';
import ExcellentExport from 'excellentexport';
import { withTranslation } from 'react-i18next';
import { TFunction } from 'i18next';

interface State {
  selectedItem: string
  loading: boolean
}

interface Props {
  t: TFunction
}

class Users extends React.Component<Props, State> {
  authProviders: { [key: string]: string } = {};
  data: User[] = [];

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
    AuthProvider.list().then(providers => {
      providers.forEach(provider => {
        this.authProviders[provider.id] = provider.name;
      });
      User.list().then(list => {
        this.data = list;
        this.setState({ loading: false });
      });
    });
  }

  onItemSelect = (user: User) => {
    this.setState({ selectedItem: user.id });
  }

  renderItem = (user: User) => {
    let authProvider = "";
    if (user.requirePassword) {
      authProvider = this.props.t("password");
    } else if (this.authProviders[user.authProviderId]) {
      authProvider = this.authProviders[user.authProviderId];
    }
    let role = this.props.t("roleUser");
    if (user.role === User.UserRoleSpaceAdmin) {
      role = this.props.t("roleSpaceAdmin");
    }
    if (user.role === User.UserRoleOrgAdmin) {
      role = this.props.t("roleOrgAdmin");
    }
    return (
      <tr key={user.id} onClick={() => this.onItemSelect(user)}>
        <td>{user.email}</td>
        <td>{role}</td>
        <td>{authProvider}</td>
      </tr>
    );
  }

  exportTable = (e: any) => {
    return ExcellentExport.convert(
      { anchor: e.target, filename: "seatsurfing-users", format: "xlsx"},
      [{name: "Seatsurfing Users", from: {table: "datatable"}}]
    );
  }

  render() {
    if (this.state.selectedItem) {
      return <Navigate replace={true} to={`/users/${this.state.selectedItem}`} />
    }
    // eslint-disable-next-line
    let downloadButton = <a download="seatsurfing-users.xlsx" href="#" className="btn btn-sm btn-outline-secondary" onClick={this.exportTable}><IconDownload className="feather" /> {this.props.t("download")}</a>;
    let buttons = (
      <>
        {this.data && this.data.length > 0 ? downloadButton : <></>}
        <Link to="/users/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> {this.props.t("add")}</Link>
      </>
    );

    if (this.state.loading) {
      return (
        <FullLayout headline={this.props.t("users")} buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline={this.props.t("users")} buttons={buttons}>
          <p>{this.props.t("noRecords")}</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline={this.props.t("users")} buttons={buttons}>
        <Table striped={true} hover={true} className="clickable-table" id="datatable">
          <thead>
            <tr>
              <th>{this.props.t("username")}</th>
              <th>{this.props.t("role")}</th>
              <th>{this.props.t("loginMeans")}</th>
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

export default withTranslation()(Users as any);
