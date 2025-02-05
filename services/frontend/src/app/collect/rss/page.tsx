'use client'
import Form from 'next/form'
import { useState } from 'react'

export default function Page() {

  const [urls, setUrls] = useState([{ id: 1 }]);

  return <div>
    <h1 className="text-4xl font-bold">RSS</h1>
    <div className="m-4"></div>
    <Form action="/collect/rss">
      <fieldset className="bg-base-100 border border-base-300 p-4 rounded-box">
        <input
          type="text"
          className="input w-full"
          placeholder='Title'
        />
      </fieldset>
      <fieldset className="bg-base-100 border border-base-300 p-4 rounded-box">
        <legend className="fieldset-legend">Feeds</legend>
        <button className="btn" type="button" onClick={() => setUrls([...urls, { id: urls.length + 1 }])}>Add URL</button>
        {urls.map((ea) => <label key={ea.id} className="floating-label my-4">
          <span>URL</span>
          <input
            type="url"
            className="input validator w-full"
            required={ea.id === 1}
            placeholder="URL"
            name={'url-' + ea.id}
            pattern="^(https?://)?([a-zA-Z0-9]([a-zA-Z0-9-].*[a-zA-Z0-9])?.)+[a-zA-Z].*$"
            title="Must be valid URL"
          />
          <p className="validator-hint">Must be valid URL</p>
        </label>)}
      </fieldset>
      <div className="m-4"></div>
      <button type="submit" className="btn">Submit</button>
    </Form>
  </div>
}