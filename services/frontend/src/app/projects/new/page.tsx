import Form from "next/form";
import { ProjectStatus } from "@/models/project";
import { newProject } from "./actions";

export default async function ProjectNewPage() {
    return <div className="container mx-auto max-w-4xl px-5">
        <h1 className="text-2xl font-bold mb-4">New Project</h1>
        <Form action={newProject} className="bg-base-200 rounded-box p-4">
            <fieldset className="fieldset flex flex-col lg:flex-row lg:gap-4">
                <div className="w-full lg:w-1/3">
                    <label className="fieldset-label">Project Number</label>
                    <input
                        type="text"
                        className="input validator w-full"
                        placeholder="Number"
                        name="number"
                        title="Number is required"
                        defaultValue='2025-002'
                        required
                    />
                    <div className="validator-hint">Required</div>
                </div>
                <div className="w-full">
                    <label className="fieldset-label">Title</label>
                    <input
                        type="text"
                        className="input validator w-full"
                        placeholder="Title"
                        name="title"
                        title="Title is required"
                        required
                    />
                    <div className="validator-hint">Required</div>
                </div>
            </fieldset>
            <fieldset className="fieldset">
                <label className="fieldset-label">Status</label>
                <select className="select w-full validator" name="status" required>
                    {Object.keys(ProjectStatus).map((key) => <option key={key} value={key}>{ProjectStatus[key].label}</option>)}
                </select>
                <div className="validator-hint">Required</div>
            </fieldset>
            <div className="flex items-center mt-4 h-12 gap-2 justify-between">
                <input type="submit" className="btn btn-primary h-full" defaultValue="Create" />
                <a href="/projects" className="btn h-full">Cancel</a>
            </div>
        </Form>
    </div>
}