import React from 'react';
import FullLayout from '../components/FullLayout';
import { Table } from 'react-bootstrap';
import { Plus as IconPlus } from 'react-feather';
import { Link, Redirect } from 'react-router-dom';
import Loading from '../components/Loading';
import { User, AuthProvider } from 'flexspace-commons';

interface State {
  selectedItem: string
  loading: boolean
}

export default class Users extends React.Component<{}, State> {
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
      authProvider = "Kennwort";
    } else if (this.authProviders[user.authProviderId]) {
      authProvider = this.authProviders[user.authProviderId];
    }
    return (
      <tr key={user.id} onClick={() => this.onItemSelect(user)}>
        <td>{user.email}</td>
        <td>{user.admin ? "Ja" : ""}</td>
        <td>{authProvider}</td>
      </tr>
    );
  }

  render() {
    if (this.state.selectedItem) {
      return <Redirect to={`/users/${this.state.selectedItem}`} />
    }

    let buttons = <Link to="/users/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> Neu</Link>;

    if (this.state.loading) {
      return (
        <FullLayout headline="Benutzer" buttons={buttons}>
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline="Benutzer" buttons={buttons}>
          <p>Keine Datensätze gefunden.</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline="Benutzer" buttons={buttons}>
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>Benutzername</th>
              <th>Admin</th>
              <th>Anmeldung</th>
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
