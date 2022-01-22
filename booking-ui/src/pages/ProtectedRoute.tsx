import './Login.css';
import { Navigate } from 'react-router-dom';
import { Ajax } from 'flexspace-commons';
import NavBar from '../components/NavBar';

export function ProtectedRoute({ children }: { children: JSX.Element }) {
  if (!Ajax.CREDENTIALS.accessToken) {
    return <Navigate to="/login" replace={true} />;
  }
  return (
    <>
      <NavBar />
      {children}
    </>
  );
}
