'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import useWebSocket from '../../../utils/useWebSocket'
import moment from 'moment'

export default function Page() {
  const { job } = useParams<{ job: string }>()
  const { messages, sendMessage } = useWebSocket('ws://localhost:3000/api/ws')
  const parsedMessages = messages.map((message) => JSON.parse(message)).sort((a, b) => b.timestamp - a.timestamp)
  return <div>
    {parsedMessages.map((message, i) => <div key={i}>
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