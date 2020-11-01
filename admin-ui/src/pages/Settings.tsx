import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { User, Organization, AuthProvider, Settings as OrgSettings, Domain, Ajax } from 'flexspace-commons';
import { Form, Col, Row, Table, Button, Alert, InputGroup, Popover, OverlayTrigger } from 'react-bootstrap';
import { Link, Redirect } from 'react-router-dom';
import { Plus as IconPlus, Save as IconSave } from 'react-feather';

interface State {
  allowAnyUser: boolean
  maxBookingsPerUser: number
  maxDaysInAdvance: number
  maxBookingDurationHours: number
  subscriptionActive: boolean
  subscriptionMaxUsers: number
  selectedAuthProvider: string
  loading: boolean
  submitting: boolean
  saved: boolean
  error: boolean
  newDomain: string
  domains: Domain[]
  userDomain: string
}

export default class Settings extends React.Component<{}, State> {
  org: Organization | null;
  authProviders: AuthProvider[];

  constructor(props: any) {
    super(props);
    this.org = null;
    this.authProviders = [];
    this.state = {
      allowAnyUser: true,
      maxBookingsPerUser: 0,
      maxBookingDurationHours: 0,
      maxDaysInAdvance: 0,
      subscriptionActive: false,
      subscriptionMaxUsers: 0,
      selectedAuthProvider: "",
      loading: true,
      submitting: false,
      saved: false,
      error: false,
      newDomain: "",
      domains: [],
      userDomain: ""
    };
  }

  componentDidMount = () => {
    this.loadSettings();
    this.loadItems();
  }

  loadItems = () => {
    User.getSelf().then(user => {
      let userDomain = user.email.substring(user.email.indexOf("@")+1).toLowerCase();
      Organization.get(user.organizationId).then(org => {
        this.org = org;
        Domain.list(org.id).then(domains => {
          this.setState({
            domains: domains,
            userDomain: userDomain
          });
          AuthProvider.list().then(list => {
            this.authProviders = list;
            this.setState({ loading: false });
          });
        });
      });
    });
  }

  loadSettings = () => {
    OrgSettings.list().then(settings => {
      let state: any = {};
      settings.forEach(s => {
        if (s.name === "allow_any_user") state.allowAnyUser = (s.value === "1");
        if (s.name === "max_bookings_per_user") state.maxBookingsPerUser = window.parseInt(s.value);
        if (s.name === "max_days_in_advance") state.maxDaysInAdvance = window.parseInt(s.value);
        if (s.name === "max_booking_duration_hours") state.maxBookingDurationHours = window.parseInt(s.value);
        if (s.name === "subscription_active") state.subscriptionActive = (s.value === "1");
        if (s.name === "subscription_max_users") state.subscriptionMaxUsers = window.parseInt(s.value);
      });
      this.setState({
        ...this.state,
        ...state
      });
    });
  }

  onSubmit = (e: any) => {
    e.preventDefault();
    this.setState({
      submitting: true,
      saved: false,
      error: false
    });
    let payload = [
      new OrgSettings("allow_any_user", this.state.allowAnyUser ? "1" : "0"),
      new OrgSettings("max_bookings_per_user", this.state.maxBookingsPerUser.toString()),
      new OrgSettings("max_days_in_advance", this.state.maxDaysInAdvance.toString()),
      new OrgSettings("max_booking_duration_hours", this.state.maxBookingDurationHours.toString())
    ];
    OrgSettings.setAll(payload).then(() => {
      this.setState({
        submitting: false,
        saved: true
      });
    }).catch(() => {
      this.setState({
        submitting: false,
        error: true
      });
    });
  }

  onAuthProviderSelect = (e: AuthProvider) => {
    this.setState({ selectedAuthProvider: e.id });
  }

  getAuthProviderTypeLabel = (providerType: number): string => {
    switch (providerType) {
      case 1: return "OAuth 2";
      default: return "Unknown";
    }
  }

  renderAuthProviderItem = (e: AuthProvider) => {
    return (
      <tr key={e.id} onClick={() => this.onAuthProviderSelect(e)}>
        <td>{e.name}</td>
        <td>{this.getAuthProviderTypeLabel(e.providerType)}</td>
      </tr>
    );
  }

  verifyDomain = (domainName: string) => {
    document.body.click();
    this.state.domains.forEach(domain => {
      if (domain.domain === domainName) {
        domain.verify().then(() => {
          Domain.list(domain.organizationId).then(domains => this.setState({ domains: domains }));
        }).catch(e => {
          alert("Fehler beim Bestätigen der Domain " + domainName + ": Bitte stellen Sie sicher, dass der notwendige DNS-TXT-Record korrekt eingerichtet ist.");
        })
      }
    });
  }

  isValidDomain = () => {
    if (this.state.newDomain.indexOf(".") < 3) {
      return false;
    }
    let lastIndex = this.state.newDomain.length - 3;
    if (lastIndex < 3) {
      lastIndex = 3;
    }
    if (this.state.newDomain.lastIndexOf(".") > lastIndex) {
      return false;
    }
    return true;
  }

  addDomain = () => {
    if (!this.isValidDomain() || !this.org) {
      return;
    }
    Domain.add(this.org.id, this.state.newDomain).then(() => {
      Domain.list(this.org ? this.org.id : "").then(domains => this.setState({ domains: domains }));
      this.setState({ newDomain: "" });
    }).catch(() => {
      alert("Fehler beim Hinzufügen der Domain.");
    });
  }

  removeDomain = (domainName: string) => {
    if (!window.confirm("Soll die Domain " + domainName + " wirklich entfernt werden?")) {
      return;
    }
    this.state.domains.forEach(domain => {
      if (domain.domain === domainName) {
        domain.delete().then(() => {
          Domain.list(this.org ? this.org.id : "").then(domains => this.setState({ domains: domains }));
        }).catch(() => alert("Fehler beim Entfernen der Domain."));
      }
    });
  }

  handleNewDomainKeyDown = (target: any) => {
    if (target.key === "Enter") {
      target.preventDefault();
      this.addDomain();
    }
  }

  deleteOrg = () => {
    if (window.confirm("Diese Organisation unwiederbringlich löschen?")) {
      if (window.confirm("Sind Sie ganz sicher? Wenn Sie diese Organisation löschen, werden alle Buchungen, Bereiche, Plätze und Benutzer unwiederbringlich gelöscht. Eine Wiederherstellung ist nicht möglich.")) {
        this.org?.delete().then(() => {
          Ajax.JWT = "";
          window.sessionStorage.removeItem("jwt");
          window.location.href = "/admin/";
        });
      }
    }
  }

  manageSubscription = () => {
    let windowRef = window.open();
    this.org?.getSubscriptionManagementURL().then(url => {
      if (windowRef) {
        windowRef.location.href = url;
      }
    }).catch(() => {
      if (windowRef) {
        windowRef?.close();
      }
      alert("Etwas ist schief gegangen. Bitte probieren Sie es später erneut.");
    });
  }

  render() {
    if (this.state.selectedAuthProvider) {
      return <Redirect to={`/settings/auth-providers/${this.state.selectedAuthProvider}`} />
    }

    if (this.state.loading) {
      return (
        <FullLayout headline="Einstellungen">
          <Loading />
        </FullLayout>
      );
    }

    let domains = this.state.domains.map(domain => {
      let verify = <></>;
      let popoverId = "popover-domain-" + domain.domain;
      const popover = (
        <Popover id={popoverId}>
          <Popover.Title as="h3">Domain bestätigen</Popover.Title>
          <Popover.Content>
            <div>Um die Domain <strong>{domain.domain}</strong> zu bestätigen, fügen Sie im DNS-Server der Domain bitte folgenden TXT-Record hinzu:</div>
            <div>&nbsp;</div>
            <div><strong>seatsurfing-verification={domain.verifyToken}</strong></div>
            <div>&nbsp;</div>
            <Button variant="primary" size="sm" onClick={() => this.verifyDomain(domain.domain)}>Jetzt bestätigen</Button>
          </Popover.Content>
        </Popover>
      );
      if (!domain.active) {
        verify = (
          <OverlayTrigger trigger="click" placement="auto" overlay={popover} rootClose={true}>
            <Button variant="primary" size="sm">Bestätigen</Button>
          </OverlayTrigger>
        );
      }
      let key = "domain-" + domain.domain;
      let canDelete = domain.domain.toLowerCase() !== this.state.userDomain;
      return (
        <Form.Group key={key}>
          {domain.domain}
            &nbsp;
          <Button variant="danger" size="sm" onClick={() => this.removeDomain(domain.domain)} disabled={!canDelete}>Entfernen</Button>
            &nbsp;
          {verify}
        </Form.Group>
      );
    });

    let authProviderRows = this.authProviders.map(item => this.renderAuthProviderItem(item));
    let authProviderTable = <p>Keine Datensätze gefunden.</p>;
    if (authProviderRows.length > 0) {
      authProviderTable = (
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Typ</th>
            </tr>
          </thead>
          <tbody>
            {authProviderRows}
          </tbody>
        </Table>
      );
    }

    let subscription = <></>;
    if (this.state.subscriptionActive) {
      subscription = (
        <>
          <p>Sie haben ein aktives Abonnement von Seatsurfing mit bis zu {this.state.subscriptionMaxUsers} Benutzern.</p>
          <p><Button variant="primary" onClick={this.manageSubscription}>Abonnement verwalten</Button></p>
        </>
      );
    } else {
      subscription = (
        <>
          <p>Sie verwenden aktuell die kostenfreie Version von Seatsurfing mit bis zu {this.state.subscriptionMaxUsers} Benutzern.</p>
          <p><Button variant="primary" onClick={this.manageSubscription}>Abonnement verwalten</Button></p>
        </>
      );
    }

    let dangerZone = (
      <>
        <Button className="btn btn-danger" onClick={this.deleteOrg}>Organisation löschen</Button>
      </>
    );

    let hint = <></>;
    if (this.state.saved) {
      hint = <Alert variant="success">Eintrag wurde aktualisiert.</Alert>
    } else if (this.state.error) {
      hint = <Alert variant="danger">Fehler beim Speichern, bitte kontrollieren Sie die Angaben.</Alert>
    }

    let buttonSave = <Button className="btn-sm" variant="outline-secondary" type="submit" form="form"><IconSave className="feather" /> Speichern</Button>;
    let contactName = "";
    if (this.org) {
      contactName = this.org.contactFirstname + " " + this.org.contactLastname + " ("+this.org.contactEmail+")";
    }

    return (
      <FullLayout headline="Einstellungen" buttons={buttonSave}>
        <Form onSubmit={this.onSubmit} id="form">
          {hint}
          <Form.Group as={Row}>
            <Form.Label column sm="2">Organisation</Form.Label>
            <Col sm="4">
              <Form.Control plaintext={true} readOnly={true} defaultValue={this.org?.name} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Primärkontakt</Form.Label>
            <Col sm="4">
              <Form.Control plaintext={true} readOnly={true} defaultValue={contactName} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Col sm="6">
              <Form.Check type="checkbox" id="check-allowAnyUser" label="Login aller authentifzierbaren Benutzer erlauben" checked={this.state.allowAnyUser} onChange={(e: any) => this.setState({ allowAnyUser: e.target.checked })} />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Buchungen je Nutzer</Form.Label>
            <Col sm="4">
              <Form.Control type="number" value={this.state.maxBookingsPerUser} onChange={(e: any) => this.setState({ maxBookingsPerUser: e.target.value })} min="1" max="9999" />
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Max. Buchungs-Vorlauf</Form.Label>
            <Col sm="4">
              <InputGroup>
                <Form.Control type="number" value={this.state.maxDaysInAdvance} onChange={(e: any) => this.setState({ maxDaysInAdvance: e.target.value })} min="0" max="9999" />
                <InputGroup.Append>
                  <InputGroup.Text>Tage</InputGroup.Text>
                </InputGroup.Append>
              </InputGroup>
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Max. Buchungs-Dauer</Form.Label>
            <Col sm="4">
              <InputGroup>
                <Form.Control type="number" value={this.state.maxBookingDurationHours} onChange={(e: any) => this.setState({ maxBookingDurationHours: e.target.value })} min="0" max="9999" />
                <InputGroup.Append>
                  <InputGroup.Text>Stunden</InputGroup.Text>
                </InputGroup.Append>
              </InputGroup>
            </Col>
          </Form.Group>
          <Form.Group as={Row}>
            <Form.Label column sm="2">Domains</Form.Label>
            <Col sm="4">
              {domains}
              <InputGroup size="sm">
                <Form.Control type="text" value={this.state.newDomain} onChange={(e: any) => this.setState({ newDomain: e.target.value })} placeholder="ihre-domain.de" onKeyDown={this.handleNewDomainKeyDown} />
                <InputGroup.Append>
                  <Button variant="outline-secondary" onClick={this.addDomain} disabled={!this.isValidDomain()}>Domain hinzufügen</Button>
                </InputGroup.Append>
              </InputGroup>
            </Col>
          </Form.Group>
        </Form>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Abonnement</h1>
        </div>
        {subscription}
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Auth Providers</h1>
          <div className="btn-toolbar mb-2 mb-md-0">
            <div className="btn-group mr-2">
              <Link to="/settings/auth-providers/add" className="btn btn-sm btn-outline-secondary"><IconPlus className="feather" /> Neu</Link>
            </div>
          </div>
        </div>
        {authProviderTable}
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
          <h1 className="h2">Kritische Funktionen</h1>
        </div>
        {dangerZone}
      </FullLayout>
    );
  }
}
