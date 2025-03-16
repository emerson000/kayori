'use server'

import MessageCard from '@/components/MessageCard'
import { getProject } from '@/services/projectService'
import ProjectHeader from '@/components/projects/projectHeader'
import { getJob, getJobArtifacts } from '@/services/jobService'
import { notFound } from 'next/navigation'

export default async function Page({ params }: { params: Promise<{ id: string, job: string }> }) {
  const { id, job } = await params;
  const project = await getProject(id);
  const jobRecord = await getJob(id, job)
  const artifacts = await getJobArtifacts(id, job)
  if (!project || !jobRecord || !artifacts) {
    notFound()
  }
  return <div>
    <ProjectHeader project={project} currentPage="collect" />
    <h2 className="text-2xl font-bold">{jobRecord.title}</h2>
    {artifacts.map((message: any, i: number) => <MessageCard key={i} message={message} />)}
  </div>
}