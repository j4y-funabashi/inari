import React from 'react';
import { timelineMonthResponse } from '../apiClient';
import { Link } from 'react-router-dom';

export interface MediaTimelineMonthProps {
  mediaTimeline: timelineMonthResponse
}

const MediaTimelineMonth: React.FunctionComponent<MediaTimelineMonthProps> = (props: MediaTimelineMonthProps) => {
  const { mediaTimeline } = props

  const header = (
    <div>
      <h1>{mediaTimeline.collection_meta.title}</h1>
      <small>{mediaTimeline.collection_meta.media_count}</small>
    </div>
  )
  const media = mediaTimeline.media.map((m) => {
    return (
      <Link to={`/media/${m.id}`}><img src={`/${m.media_src.small}`} /></Link>
    )
  })

  return (
    <div>
      {header}
      <div>{media}</div>
    </div>
  )
}

export default MediaTimelineMonth
