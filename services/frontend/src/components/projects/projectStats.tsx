import { IProject } from "@/models/project";
export default function ProjectStats({ project }: { project: IProject }) {
    return <div className="w-full overflow-x-auto pt-2 pb-1 px-1">
        <div className="stats shadow outline outline-base-300 flex flex-row">
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">Critical Alerts</div>
                <div className="stat-value text-error">1</div>
            </a>
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">Alerts</div>
                <div className="stat-value text-warning">15</div>
            </a>
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">Saved Artifacts</div>
                <div className="stat-value">20</div>
            </a>
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">People</div>
                <div className="stat-value">48</div>
            </a>
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">Organizations</div>
                <div className="stat-value">12</div>
            </a>
            <a className="stat basis-40 hover:bg-base-200 active:bg-base-300" href="#">
                <div className="stat-title">News Articles</div>
                <div className="stat-value">80,100</div>
            </a>
        </div>
    </div>
}