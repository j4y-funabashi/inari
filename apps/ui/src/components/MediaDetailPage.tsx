import React from 'react';
import { useParams } from 'react-router-dom';
import { fetchMediaDetail, mediaDetailResponse } from '../apiClient';

type urlParams = {
  mediaid: string
}

const MediaDetailPage: React.FunctionComponent = () => {
  const [mediaDetailData, setMediaDetailData] = React.useState<mediaDetailResponse>({ media: { id: "", date: "", media_src: { small: "", medium: "", large: "" } } });

  const { mediaid } = useParams<urlParams>()
  console.log(mediaid)

  React.useEffect(() => {
    (async () => {
      const mediaDetailResponse = await fetchMediaDetail(mediaid)
      setMediaDetailData(mediaDetailResponse)
    })()
  }, [setMediaDetailData])

  console.log(mediaDetailData)

  return (

    <article className="vh-100 dt w-100 bg-dark-gray">
      <div className="dtc v-mid tc">
        <img src={`https://photos-dev.funabashi.co.uk/${mediaDetailData.media.media_src.large}`} alt="" />
      </div>
    </article>

  )
}

export default MediaDetailPage
