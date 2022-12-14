import React from "react";
import { useParams } from "react-router-dom";
import { MediaDetailQuery, mediaDetailResponse } from "../apiClient";
import MediaDetail from "./MediaDetail";

type urlParams = {
	mediaid: string;
};

interface MediaDetailPageProps {
	fetchMediaDetail: MediaDetailQuery;
	setMediaDetailData: React.Dispatch<React.SetStateAction<mediaDetailResponse>>;
	media: mediaDetailResponse;
}

const MediaDetailPage: React.FunctionComponent<
	React.PropsWithChildren<MediaDetailPageProps>
> = (props: MediaDetailPageProps) => {
	const { media, setMediaDetailData, fetchMediaDetail } = props;

	const { mediaid } = useParams<urlParams>();
	console.log(mediaid);

	React.useEffect(() => {
		(async () => {
			const mediaDetailResponse = await fetchMediaDetail(mediaid!);
			setMediaDetailData(mediaDetailResponse);
		})();
	}, [setMediaDetailData, mediaid, fetchMediaDetail]);

	return <MediaDetail media={media.media} />;
};

export default MediaDetailPage;
