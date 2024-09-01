import { Link } from 'react-router-dom';
import './NotFoundPage.scss';

const NotFoundPage = () => (
  <div className="not-found-page">
    <h2 className="title">Couldn't find desired page :( </h2>
    <Link className="go-home-btn" to="/">
      Go Home
    </Link>
  </div>
);

export default NotFoundPage;
