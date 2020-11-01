import React from 'react';
import './SideBar.css';
import { Home as IconHome, Users as IconUsers, Map as IconMap, Book as IconBook, Settings as SettingsIcon } from 'react-feather';
import { NavLink } from 'react-router-dom';

export default class SideBar extends React.Component {
    render() {
        return (
            <nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div className="sidebar-sticky pt-3">
                    <ul className="nav flex-column">
                        <li className="nav-item">
                            <NavLink to="/dashboard" className="nav-link" activeClassName="active"><IconHome className="feather" /> Dashboard</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/locations" className="nav-link" activeClassName="active"><IconMap className="feather" /> Bereiche</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/users" className="nav-link" activeClassName="active"><IconUsers className="feather" /> Benutzer</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/bookings" className="nav-link" activeClassName="active"><IconBook className="feather" /> Buchungen</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink to="/settings" className="nav-link" activeClassName="active"><SettingsIcon className="feather" /> Einstellungen</NavLink>
                        </li>
                    </ul>
                </div>
            </nav>
        );
    }
}
