import React from 'react';
import { collectionsResponse } from '../apiClient';
import { Link } from 'react-router-dom';

export interface MediaTimelineProps {
  mediaTimeline: collectionsResponse
  collection_type: string
}

const MediaTimeline: React.FunctionComponent<React.PropsWithChildren<MediaTimelineProps>> = (props: MediaTimelineProps) => {
  const { mediaTimeline, collection_type } = props

  const mediaMonths = mediaTimeline.map((m) => {
    const url = "/collection/" + collection_type + "/" + m.id
    return (
      <li key={m.id}><Link to={url}>{m.title}</Link> <small>({m.media_count})</small></li>
    )
  })

  return (<ol>{mediaMonths}</ol>)
}

export default MediaTimeline
