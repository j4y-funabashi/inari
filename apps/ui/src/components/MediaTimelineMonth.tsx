import React from 'react';
import { Link } from 'react-router-dom';
import { timelineMonthResponse } from '../apiClient';
import { format } from 'date-fns'

export interface MediaTimelineMonthProps {
  mediaTimeline: timelineMonthResponse
}

const MediaTimelineMonth: React.FunctionComponent<MediaTimelineMonthProps> = (props: MediaTimelineMonthProps) => {
  const { mediaTimeline } = props

  // sort by date
  mediaTimeline.media.sort(
    (a, b) => { if (a.date < b.date) { return -1 } return 1 }
  )
  // group media to days
  let dayCollections = new Map<string, timelineMonthResponse>();
  mediaTimeline.media.forEach((m) => {
    const dat = new Date(m.date)
    const datKey = format(dat, 'yyyy-MM-dd')
    const datTitle = format(dat, 'eee, do MMM')
    const dayCollection = dayCollections.get(datKey)
    if (dayCollection) {
      dayCollection.media.push(m)
      dayCollection.collection_meta.media_count = dayCollection.media.length
      dayCollections.set(datKey, dayCollection)
    } else {
      const c: timelineMonthResponse = {
        media: [m],
        collection_meta: {
          id: datKey,
          title: datTitle,
          type: "timeline_day",
          media_count: 1
        }
      }
      dayCollections.set(datKey, c)
    }
  })

  // render
  const media = Array.from(dayCollections.values()).map((v) => {
    const thumbs = v.media.map((m) => {
      return (
        <Link key={m.id} to={`/media/${m.id}`}><img src={`/${m.media_src.small}`} /></Link>
      )
    })
    return (
      <div key={v.collection_meta.id}>
        <h2>{v.collection_meta.title}</h2>
        <div>{thumbs}</div>
      </div>
    )
  })

  const header = (
    <div>
      <h1>{mediaTimeline.collection_meta.title}</h1>
      <small>{mediaTimeline.collection_meta.media_count}</small>
    </div>
  )

  return (
    <div>
      {header}
      {media}
    </div>
  )
}

export default MediaTimelineMonth
