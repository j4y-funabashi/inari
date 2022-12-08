import format from "date-fns/format";
import React from "react";
import { useParams } from "react-router-dom";
import { MediaDetailQuery, mediaDetailResponse } from "../apiClient";

type urlParams = {
	mediaid: string;
};

interface MediaDetailPageProps {
	fetchMediaDetail: MediaDetailQuery;
	setMediaDetailData: React.Dispatch<React.SetStateAction<mediaDetailResponse>>;
	media: mediaDetailResponse;
}

const MediaDetailPage: React.FunctionComponent<MediaDetailPageProps> = (
	props: MediaDetailPageProps,
) => {
	const { media, setMediaDetailData, fetchMediaDetail } = props;

	const { mediaid } = useParams<urlParams>();
	console.log(mediaid);

	React.useEffect(() => {
		(async () => {
			const mediaDetailResponse = await fetchMediaDetail(mediaid);
			setMediaDetailData(mediaDetailResponse);
		})();
	}, [setMediaDetailData, mediaid, fetchMediaDetail]);

	console.log(media);
	const dat = new Date(media.media.date);
	console.log(dat);
	const datKey = format(dat, "eee, do MMM yyyy - HH:mm");
	const location = `${media.media.location.locality}, ${media.media.location.region}`;
	const caption = media.media.caption;

	return (
		<article>
			<div>
				<img src={`${media.media.media_src.large}`} alt="" />
				<p>{caption}</p>
				<p>{datKey}</p>
				<p>{location}</p>
			</div>
			<div>
				<button type="submit">Add Caption</button>
				<button type="submit">Add Location</button>
				<button type="submit">Delete</button>
			</div>
		</article>
	);
};

export default MediaDetailPage;
