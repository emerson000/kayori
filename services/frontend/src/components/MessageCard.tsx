import moment from 'moment'

export default function MessageCard({ message }: { message: any }) {
  return (
    <div className="card card-dash bg-base-200 w-full my-6">
      <div className="card-body">
        <h2 className="card-title">{message.title}</h2>
        <p>{moment(message.timestamp).fromNow()}</p>
        <div dangerouslySetInnerHTML={{ __html: message.description }}></div>
        <div className="card-actions justify-end">
          <a className="btn btn-primary" target="_blank" href={message.url}>View</a>
        </div>
      </div>
    </div>
  )
}
