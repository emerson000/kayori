'use client'
import Form from "next/form"
import { useState } from "react"
import { createRssFeed } from "./actions"

export default function RssForm({ id }: { id: string }) {
    const [urls, setUrls] = useState([{ id: 1 }]);
    const [showSchedule, setShowSchedule] = useState(false);
    return (<div><Form action={(formData) => createRssFeed(formData, id)} >
        <fieldset className="bg-base-100 border border-base-300 p-4 rounded-box" >
            <input
                type="text"
                className="input validator w-full"
                placeholder='Title'
                name="title"
                title="Title is required"
                required
            />
            <p className="validator-hint" > Title is required </p>
        </fieldset>
        < fieldset className="bg-base-100 border border-base-300 p-4 rounded-box" >
            <legend className="fieldset-legend" > Feeds </legend>
            < button className="btn" type="button" onClick={() => setUrls([...urls, { id: urls.length + 1 }])
            }> Add URL </button>
            {
                urls.map((ea) => <label key={ea.id} className="floating-label my-4" >
                    <span>URL </span>
                    < input
                        type="url"
                        className="input validator w-full"
                        required={ea.id === 1}
                        placeholder="URL"
                        name="urls[]"
                        pattern="^(https?://)?([a-zA-Z0-9]([a-zA-Z0-9-].*[a-zA-Z0-9])?.)+[a-zA-Z].*$"
                        title="Must be valid URL"
                    />
                    <p className="validator-hint" > Must be valid URL </p>
                </label>)
            }
        </fieldset>
        < fieldset className="bg-base-100 border border-base-300 p-4 rounded-box" >
            <legend className="fieldset-legend" > Schedule </legend>
            < input name="schedule" type="checkbox" value="true" checked={showSchedule} className="toggle" onChange={() => setShowSchedule(!showSchedule)} />
            {
                showSchedule && <div className="flex space-x-4" >
                    <label className="floating-label my-4 flex-1" >
                        <span>Duration </span>
                        < input type="number"
                            className="input validator"
                            placeholder="Duration"
                            name="duration"
                            required />
                    </label>
                    < label className="floating-label my-4 flex-1" >
                        <span>Interval </span>
                        < select defaultValue=""
                            className="select validator"
                            name="interval" required >
                            <option disabled value="" > Select an interval </option>
                            < option value="minutes" > Minutes </option>
                            < option value="hours" > Hours </option>
                            < option value="days" > Days </option>
                        </select>
                        < p className="validator-hint" > Required </p>
                    </label>
                </div>
            }
        </fieldset>
        < div className="m-4" > </div>
        < button type="submit" className="btn" > Submit </button>
    </Form>
    </div>);
}