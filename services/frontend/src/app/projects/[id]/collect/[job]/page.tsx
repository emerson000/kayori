'use server'

import { getJobArtifacts } from '@/utils/shared'
import MessageCard from '@/components/MessageCard'
import moment from 'moment'
import { getProject } from '@/services/projectService'
import { IProject } from '@/models/project'
import ProjectHeader from '@/components/projects/projectHeader'

export default async function Page({ params }: { params: Promise<{ id: string, job: string }> }) {
  const { id, job } = await params;
  const project = await getProject(id);
  const artifacts = await getJobArtifacts(job);

  return <div>
    {project && <ProjectHeader project={project} />}
    {artifacts.map((message: any, i: number) => <MessageCard key={i} message={message} />)}
  </div>
}