import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus } from 'react-feather';
import { Link, Redirect } from 'react-router-dom';
import Loading from '../components/Loading';
import { User, AuthProvider } from 'flexspace-commons';
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
    return (
      <tr key={user.id} onClick={() => this.onItemSelect(user)}>
        <td>{user.email}</td>
        <td>{user.admin ? this.props.t("yes") : ""}</td>
        <td>{authProvider}</td>
      </tr>
    );
  }

  render() {
    if (this.state.selectedItem) {
      return <Redirect to={`/users/${this.state.selectedItem}`} />
    }

    let buttons = <Link to="/users/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> {this.props.t("add")}</Link>;

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
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>{this.props.t("username")}</th>
              <th>{this.props.t("admin")}</th>
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
