import React, { useEffect, useState } from 'react';
import { HttpError } from '../api/httpClient';
import { usersApi } from '../api/users';
import { useAuthContext } from '../context/AuthContext';

export const ProfilePage: React.FC = () => {
  const { user, token, refreshUser } = useAuthContext();
  const [bio, setBio] = useState(user?.bio ?? '');
  const [name, setName] = useState(user?.name ?? '');
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    setBio(user?.bio ?? '');
    setName(user?.name ?? '');
  }, [user]);

  if (!user || !token) {
    return null;
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setMessage(null);
    setIsSaving(true);

    try {
      await usersApi.updateProfile(token, { name, bio });
      await refreshUser();
      setMessage('Профиль обновлён');
    } catch (err) {
      if (err instanceof HttpError) {
        setError(err.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Не удалось обновить профиль');
      }
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <section className="profile-page">
      <h1>Профиль</h1>
      <form className="profile-form" onSubmit={handleSubmit}>
        <label>
          Email
          <input type="email" value={user.email} disabled />
        </label>
        <label>
          Имя
          <input
            type="text"
            value={name}
            onChange={(event) => setName(event.target.value)}
            required
          />
        </label>
        <label>
          О себе
          <textarea
            value={bio}
            onChange={(event) => setBio(event.target.value)}
            rows={4}
            placeholder="Расскажите о себе"
          />
        </label>
        {error ? <p className="form-error">{error}</p> : null}
        {message ? <p className="form-success">{message}</p> : null}
        <button type="submit" disabled={isSaving}>
          {isSaving ? 'Сохраняем...' : 'Сохранить изменения'}
        </button>
      </form>
    </section>
  );
};

export default ProfilePage;
