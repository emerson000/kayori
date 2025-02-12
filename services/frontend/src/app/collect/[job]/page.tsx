'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import useWebSocket from '../../../utils/useWebSocket'
import { getJobArtifacts } from '../../../utils/shared'
import MessageCard from '../../../components/MessageCard'

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
    {combinedMessages.map((message: any, i: number) => <MessageCard key={i} message={message} />)}
  </div>
}