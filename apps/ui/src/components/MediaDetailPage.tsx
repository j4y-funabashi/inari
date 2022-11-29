import format from "date-fns/format";
import React from "react";
import { useParams } from "react-router-dom";
import { MediaDetailQuery, mediaDetailResponse } from "../apiClient";

type urlParams = {
	mediaid: string;
};

interface MediaDetailPageProps {
	fetchMediaDetail: MediaDetailQuery;
}

const MediaDetailPage: React.FunctionComponent<MediaDetailPageProps> = (
	props: MediaDetailPageProps,
) => {
	const [media, setMediaDetailData] = React.useState<mediaDetailResponse>({
		media: {
			id: "",
			date: "1984-01-28T11:00:00",
			media_src: { small: "", medium: "", large: "" },
		},
	});

	const { mediaid } = useParams<urlParams>();
	console.log(mediaid);

	React.useEffect(() => {
		(async () => {
			const mediaDetailResponse = await props.fetchMediaDetail(mediaid);
			setMediaDetailData(mediaDetailResponse);
		})();
	}, [setMediaDetailData, mediaid, props]);

	console.log(media);
	const dat = new Date(media.media.date);
	console.log(dat);
	const datKey = format(dat, "eee, do MMM yyyy - HH:mm");

	return (
		<article>
			<div>
				<img src={`${media.media.media_src.large}`} alt="" />
				<p>{datKey}</p>
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
