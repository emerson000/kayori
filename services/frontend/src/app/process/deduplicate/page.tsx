'use client';
import Form from "next/form";
import { useState } from "react";
import { createDeduplicateTask } from "./action";

export default function Page() {
    const [jobs, setJobs] = useState<number[]>([1]);
    return <div>
        <h1 className="text-4xl font-bold">Deduplicate</h1>
        <div className="m-4"></div>
        <Form action={createDeduplicateTask}>
            <fieldset className="bg-base-100 border border-base-300 p-4 rounded-box">
                <legend className="fieldset-legend">Settings</legend>
                <label className="floating-label my-4">
                    <span>Field to Compare</span>
                    <input type="text"
                        className="input validator w-full"
                        placeholder='Field To Compare'
                        name="field"
                        title="Field to compare is required"
                        defaultValue="url"
                        required
                    />
                    <p className="validator-hint">Field to compare is required</p>
                </label>
            </fieldset>
            <fieldset className="bg-base-100 border border-base-300 p-4 rounded-box">
                <legend className="fieldset-legend">Jobs</legend>
                <button className="btn" type="button" onClick={() => setJobs([...jobs, jobs.length + 1])}>Add Job</button>
                {jobs.map((ea) => <label key={ea} className="floating-label my-4">
                    <span>Job ID</span>
                    <input
                        type="text"
                        className="input validator w-full"
                        required={ea === 1}
                        placeholder="Job ID"
                        name="jobs[]"
                        title="Enter a job ID"
                    />
                    <p className="validator-hint">Enter a job ID</p>
                </label>)}
            </fieldset>
            <div className="m-4"></div>
            <button type="submit" className="btn">Submit</button>
        </Form>
    </div>
}