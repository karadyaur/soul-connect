import React, { useEffect, useState } from 'react';
import { postsApi, type Post } from '../api/posts';
import { HttpError } from '../api/httpClient';
import { useAuthContext } from '../context/AuthContext';

export const FeedPage: React.FC = () => {
  const { token, user } = useAuthContext();
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newPost, setNewPost] = useState('');
  const [isPublishing, setIsPublishing] = useState(false);

  useEffect(() => {
    if (!token) {
      return;
    }

    const fetchPosts = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const data = await postsApi.getFeed(token);
        setPosts(data);
      } catch (err) {
        if (err instanceof HttpError) {
          setError(err.message);
        } else if (err instanceof Error) {
          setError(err.message);
        } else {
          setError('Не удалось загрузить ленту');
        }
      } finally {
        setIsLoading(false);
      }
    };

    void fetchPosts();
  }, [token]);

  const handleCreatePost = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token || !newPost.trim()) {
      return;
    }

    setIsPublishing(true);
    try {
      const created = await postsApi.create(token, { content: newPost.trim() });
      setPosts((prev) => [created, ...prev]);
      setNewPost('');
    } catch (err) {
      if (err instanceof HttpError) {
        setError(err.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Не удалось опубликовать пост');
      }
    } finally {
      setIsPublishing(false);
    }
  };

  return (
    <section className="feed-page">
      <header className="feed-header">
        <h1>Лента</h1>
        {user ? <p>Добро пожаловать, {user.name}!</p> : null}
      </header>

      <form className="feed-form" onSubmit={handleCreatePost}>
        <textarea
          placeholder="Поделитесь своими мыслями..."
          value={newPost}
          onChange={(event) => setNewPost(event.target.value)}
          rows={3}
          disabled={!token || isPublishing}
        />
        <button type="submit" disabled={!newPost.trim() || isPublishing}>
          {isPublishing ? 'Публикация...' : 'Опубликовать'}
        </button>
      </form>

      {isLoading ? <p>Загрузка ленты...</p> : null}
      {error ? <p className="form-error">{error}</p> : null}

      <ul className="feed-list">
        {posts.map((post) => (
          <li key={post.id} className="feed-item">
            <header>
              <strong>{post.authorName}</strong>
              <time dateTime={post.createdAt}>
                {new Date(post.createdAt).toLocaleString()}
              </time>
            </header>
            <p>{post.content}</p>
          </li>
        ))}
      </ul>
    </section>
  );
};

export default FeedPage;
