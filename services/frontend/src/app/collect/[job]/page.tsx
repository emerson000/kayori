'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import useWebSocket from '../../../utils/useWebSocket'
import { getJobArtifacts } from '../../../utils/shared'
import moment from 'moment'

export default function Page() {
  const { job } = useParams<{ job: string }>()
  const { messages, sendMessage } = useWebSocket('ws://localhost:3000/api/ws')
  const [combinedMessages, setMessages] = useState(messages);

  useEffect(() => {
    async function fetchArtifacts() {
      const artifacts = await getJobArtifacts(job);
      const parsedMessages = messages.map((message) => JSON.parse(message))
      const combinedMessages = [...(Array.isArray(parsedMessages) ? parsedMessages : []), ...(Array.isArray(artifacts) ? artifacts : [])].reduce((acc, current) => {
        const x = acc.find(item => item.artifact_id === current.artifact_id);
        if (!x) {
          return acc.concat([current]);
        } else {
          return acc;
        }
      }, []);
      const sortedMessages = combinedMessages.sort((a: any, b: any) => b.timestamp - a.timestamp);
      setMessages(sortedMessages);
    }
    fetchArtifacts();
  }, [job, messages]);

  return <div>
    {combinedMessages.map((message: any, i: number) => <div key={i}>
      <div className="card card-dash bg-base-200 w-full my-6">
        <div className="card-body">
          <h2 className="card-title">{message.title}</h2>
          <p>{moment.unix(message.timestamp / 1000).fromNow()}</p>
          <p dangerouslySetInnerHTML={{ __html: message.description }}></p>
          <div className="card-actions justify-end">
            <a className="btn btn-primary" target="_blank" href={message.url}>View</a>
          </div>
        </div>
      </div>
    </div>)}
  </div>
}