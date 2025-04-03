import { notFound } from 'next/navigation';
import { getProject } from '@/services/projectService';
import ProjectHeader from '@/components/projects/projectHeader';
import RssForm from './rssForm';
import { IProject } from '@/models/project';
export default async function Page({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const project = await getProject(id);
  if (!project) {
    notFound();
  }

  return <div>
    <ProjectHeader project={project as IProject} currentPage="collect" />
    <h1 className="text-4xl font-bold">RSS</h1>
    <div className="m-4">
      <RssForm id={id} />
    </div>
  </div>
}