import { useNavigate } from 'react-router-dom';

export function withNavigate(Children) {
  return (props) => {
    const navigate = useNavigate();
    return <Children {...props} navigate={navigate} />
  }
};
