import { ProjectStatus, Project } from "@/models/project";
import Form from "next/form";
import DeleteProjectButton from "./deleteProjectButton";

interface ProjectFormProps {
    action: (formData: FormData) => Promise<void>;
    project?: Project;
}

export default function ProjectForm({ action, project }: ProjectFormProps) {
    return <Form action={action} className="bg-base-200 rounded-box p-4">
        <input type="hidden" name="id" value={project?.id} />
        <fieldset className="fieldset flex flex-col lg:flex-row lg:gap-4">
            <div className="w-full lg:w-1/3">
                <label className="fieldset-label">Project Number</label>
                <input
                    type="text"
                    className="input validator w-full"
                    placeholder="Number"
                    name="number"
                    title="Number is required"
                    defaultValue={project?.number}
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
                    defaultValue={project?.title}
                    required
                />
                <div className="validator-hint">Required</div>
            </div>
        </fieldset>
        <fieldset className="fieldset">
            <label className="fieldset-label">Status</label>
            <select className="select w-full validator" name="status" required defaultValue={project?.status}>
                {Object.keys(ProjectStatus).map((key) => <option key={key} value={key}>{ProjectStatus[key].label}</option>)}
            </select>
            <div className="validator-hint">Required</div>
        </fieldset>
        <div className="flex items-center mt-4 h-12 gap-2 justify-between">
            <input type="submit" className="btn btn-primary" defaultValue={project ? "Save" : "Create"} />
            <div className="flex items-center gap-2 h-full">
                {project && <DeleteProjectButton project={{ id: project.id, title: project.title, status: project.status }} />}
                <a href="/projects" className="btn">Cancel</a>
            </div>
        </div>
    </Form>
}