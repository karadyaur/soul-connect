import React from 'react';
import { Link } from 'react-router-dom';
import { useAuthContext } from '../context/AuthContext';

export const Layout: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { user, logout } = useAuthContext();

  return (
    <div className="app-shell">
      <header className="app-header">
        <nav className="app-nav">
          <Link to="/" className="logo">
            SoulConnect
          </Link>
          <div className="app-nav__links">
            {user ? (
              <>
                <Link to="/">Лента</Link>
                <Link to="/profile">Профиль</Link>
                <button type="button" onClick={logout} className="link-button">
                  Выйти
                </button>
              </>
            ) : (
              <>
                <Link to="/login">Войти</Link>
                <Link to="/register">Регистрация</Link>
              </>
            )}
          </div>
        </nav>
      </header>
      <main className="app-content">{children}</main>
    </div>
  );
};

export default Layout;
