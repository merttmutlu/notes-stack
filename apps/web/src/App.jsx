import { useEffect, useState } from "react";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || "";

const emptyForm = {
  title: "",
  content: "",
};

export default function App() {
  const [notes, setNotes] = useState([]);
  const [form, setForm] = useState(emptyForm);
  const [editingId, setEditingId] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    void loadNotes();
  }, []);

  async function loadNotes() {
    setLoading(true);
    setError("");

    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/notes`);
      if (!response.ok) {
        throw new Error("failed to fetch notes");
      }

      const payload = await response.json();
      setNotes(payload.data || []);
    } catch (fetchError) {
      setError(fetchError.message || "failed to fetch notes");
    } finally {
      setLoading(false);
    }
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setSubmitting(true);
    setError("");

    const method = editingId ? "PUT" : "POST";
    const url = editingId
      ? `${apiBaseUrl}/api/v1/notes/${editingId}`
      : `${apiBaseUrl}/api/v1/notes`;

    try {
      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(form),
      });

      if (!response.ok) {
        const payload = await response.json().catch(() => null);
        throw new Error(payload?.error || "request failed");
      }

      setForm(emptyForm);
      setEditingId("");
      await loadNotes();
    } catch (submitError) {
      setError(submitError.message || "request failed");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id) {
    setError("");

    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/notes/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        const payload = await response.json().catch(() => null);
        throw new Error(payload?.error || "failed to delete note");
      }

      if (editingId === id) {
        setEditingId("");
        setForm(emptyForm);
      }

      await loadNotes();
    } catch (deleteError) {
      setError(deleteError.message || "failed to delete note");
    }
  }

  function startEdit(note) {
    setEditingId(note.id);
    setForm({
      title: note.title,
      content: note.content,
    });
  }

  function resetForm() {
    setEditingId("");
    setForm(emptyForm);
    setError("");
  }

  return (
    <main className="page">
      <section className="hero">
        <p className="eyebrow">notes-stack</p>
        <h1>Notes Playground</h1>
        <p className="lede">
          A small React client for exercising the Go API and Postgres-backed note
          workflow.
        </p>
      </section>

      <section className="grid">
        <form className="card form-card" onSubmit={handleSubmit}>
          <div className="card-header">
            <h2>{editingId ? "Edit note" : "Create note"}</h2>
            {editingId ? (
              <button className="ghost-button" type="button" onClick={resetForm}>
                Cancel
              </button>
            ) : null}
          </div>

          <label className="field">
            <span>Title</span>
            <input
              value={form.title}
              onChange={(event) =>
                setForm((current) => ({ ...current, title: event.target.value }))
              }
              placeholder="Incident checklist"
              required
            />
          </label>

          <label className="field">
            <span>Content</span>
            <textarea
              value={form.content}
              onChange={(event) =>
                setForm((current) => ({ ...current, content: event.target.value }))
              }
              placeholder="Write the note body here"
              rows="8"
              required
            />
          </label>

          <button className="primary-button" type="submit" disabled={submitting}>
            {submitting ? "Saving..." : editingId ? "Update note" : "Create note"}
          </button>

          {error ? <p className="message error">{error}</p> : null}
        </form>

        <section className="card list-card">
          <div className="card-header">
            <h2>Notes</h2>
            <button className="ghost-button" type="button" onClick={loadNotes}>
              Refresh
            </button>
          </div>

          {loading ? <p className="message">Loading notes...</p> : null}

          {!loading && notes.length === 0 ? (
            <p className="message">No notes yet. Create the first one.</p>
          ) : null}

          <ul className="notes-list">
            {notes.map((note) => (
              <li className="note-item" key={note.id}>
                <div className="note-meta">
                  <h3>{note.title}</h3>
                  <p>{note.content}</p>
                  <small>Updated {new Date(note.updated_at).toLocaleString()}</small>
                </div>

                <div className="note-actions">
                  <button
                    className="ghost-button"
                    type="button"
                    onClick={() => startEdit(note)}
                  >
                    Edit
                  </button>
                  <button
                    className="danger-button"
                    type="button"
                    onClick={() => handleDelete(note.id)}
                  >
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </section>
      </section>
    </main>
  );
}
