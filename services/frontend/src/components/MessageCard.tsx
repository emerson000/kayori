import moment from 'moment'
import { NewsArticle } from '../models/newsArticle'
import { JSX } from 'react'

export default function MessageCard({ message }: { message: any }) {
  return (
    <div className="card card-dash bg-base-200 w-full my-6">
      <div className="card-body">
        <div>
          <h2 className="card-title">{message.title}</h2>
          <div className="text-sm text-gray-400">
            {moment(message.timestamp).fromNow()}
          </div>
        </div>
        <div dangerouslySetInnerHTML={{ __html: message.description }}></div>
        <div className="flex flex-nowrap gap-2 overflow-x-auto py-2">
          {getSources(message)}
        </div>
      </div>
    </div>
  )
}

function getSources(message: any): JSX.Element[] {
  if (!message) return []
  if (message.entity_type !== 'news_article') return []
  if (!(message instanceof NewsArticle)) {
    message = new NewsArticle(message)
  }
  const sources = message.cluster_articles?.map((article: any) => {
    return (
      <a
        key={article.id}
        className={`badge ${message.id === article.id ? 'badge-accent' : 'badge-primary'} hover:bg-base-300 hover:border-base-900`}
        title={article.title}
        href={article.url}
        target="_blank"
      >
        {article.getSecondLevelDomain()}
      </a>
    )
  })
  const defaultSource = [<a
    key={message.id}
    className="badge badge-primary hover:bg-base-300 hover:border-base-900"
    href={message.url}
    title={message.title}
    target="_blank">
    {message.getSecondLevelDomain()}
  </a>];
  return sources && sources.length > 0 ? sources : defaultSource;
}