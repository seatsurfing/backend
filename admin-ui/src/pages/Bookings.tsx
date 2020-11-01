import React from 'react';
import FullLayout from '../components/FullLayout';
import Loading from '../components/Loading';
import { Booking, Formatting } from 'flexspace-commons';
import { Table, Form, Col, Row, Button } from 'react-bootstrap';
import { Search as IconSearch } from 'react-feather';

interface State {
  loading: boolean
  start: string
  end: string
}

export default class Bookings extends React.Component<{}, State> {
  data: Booking[];

  constructor(props: any) {
    super(props);
    this.data = [];
    this.state = {
      loading: true,
      start: Formatting.getISO8601(new Date()),
      end: Formatting.getISO8601(new Date())
    };
  }

  componentDidMount = () => {
    this.loadItems();
  }

  loadItems = () => {
    let end = new Date(this.state.end);
    end.setUTCHours(23, 59, 59);
    Booking.listFiltered(new Date(this.state.start), end).then(list => {
      this.data = list;
      this.setState({ loading: false });
    });
  }

  renderItem = (booking: Booking) => {
    return (
      <tr key={booking.id}>
        <td>{booking.user.email}</td>
        <td>{booking.space.location.name}</td>
        <td>{booking.space.name}</td>
        <td>{Formatting.getFormatterShort().format(booking.enter)}</td>
        <td>{Formatting.getFormatterShort().format(booking.leave)}</td>
      </tr>
    );
  }

  onFilterSubmit = (e: any) => {
    e.preventDefault();
    this.setState({ loading: true });
    this.loadItems();
  }

  render() {
    let buttonSearch = <Button className="btn-sm" variant="outline-secondary" type="submit" form="form"><IconSearch className="feather" /> Suchen</Button>;
    let form = (
      <Form onSubmit={this.onFilterSubmit} id="form">
        <Form.Group as={Row}>
          <Form.Label column sm="2">Beginn</Form.Label>
          <Col sm="4">
            <Form.Control type="date" value={this.state.start} onChange={(e: any) => this.setState({ start: e.target.value })} required={true} />
          </Col>
        </Form.Group>
        <Form.Group as={Row}>
          <Form.Label column sm="2">Ende</Form.Label>
          <Col sm="4">
            <Form.Control type="date" value={this.state.end} onChange={(e: any) => this.setState({ end: e.target.value })} required={true} />
          </Col>
        </Form.Group>
      </Form>
    );

    if (this.state.loading) {
      return (
        <FullLayout headline="Buchungen" buttons={buttonSearch}>
          {form}
          <Loading />
        </FullLayout>
      );
    }

    let rows = this.data.map(item => this.renderItem(item));
    if (rows.length === 0) {
      return (
        <FullLayout headline="Buchungen" buttons={buttonSearch}>
          {form}
          <p>Keine Datensätze gefunden.</p>
        </FullLayout>
      );
    }
    return (
      <FullLayout headline="Buchungen" buttons={buttonSearch}>
        {form}
        <Table striped={true} hover={true} className="clickable-table">
          <thead>
            <tr>
              <th>Benutzer</th>
              <th>Bereich</th>
              <th>Platz</th>
              <th>Beginn</th>
              <th>Ende</th>
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
