import format from "date-fns/format";
import React from "react";
import { media, mediaDetailResponse } from "../apiClient";

interface MediaDetailProps {
	media: media;
}

const MediaDetail: React.FunctionComponent<
	React.PropsWithChildren<MediaDetailProps>
> = (props: MediaDetailProps) => {
	const { media } = props;

	const dat = new Date(media.date);
	const datKey = format(dat, "eee, do MMM yyyy - HH:mm");
	const location = `${media.location.locality}, ${media.location.region}`;
	const caption = media.caption;

	return (
		<article>
			<div>
				<img src={`${media.media_src.large}`} alt="" />
				<p>{caption}</p>
				<p>{datKey}</p>
				<p>{location}</p>
			</div>
			<div>
				<button type="submit">Add Caption</button>
				<button type="submit">Add Location</button>
				<DeleteMediaButton mediaID={media.id} />
			</div>
		</article>
	);
};

interface DeleteMediaButtonProps {
	mediaID: string;
}
const DeleteMediaButton: React.FunctionComponent<
	React.PropsWithChildren<DeleteMediaButtonProps>
> = (props: DeleteMediaButtonProps) => {
	const handleDeleteMedia = () => {
		console.log(props.mediaID);
	};
	return <button onClick={handleDeleteMedia}>Delete</button>;
};

export default MediaDetail;
