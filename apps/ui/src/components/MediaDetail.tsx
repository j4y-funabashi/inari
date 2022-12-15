import format from "date-fns/format";
import React from "react";
import { media } from "../apiClient";

interface MediaDetailProps {
	media: media;
	handleDelete: () => void;
	handleClose: () => void;
}

const MediaDetail: React.FunctionComponent<
	React.PropsWithChildren<MediaDetailProps>
> = (props: MediaDetailProps) => {
	const { media, handleDelete, handleClose } = props;

	const dat = new Date(media.date);
	const datKey = format(dat, "eee, do MMM yyyy - HH:mm");
	const location = `${media.location.locality}, ${media.location.region}`;
	const caption = media.caption;

	return (
		<article>
			<button type="submit">Prev</button>
			<button onClick={handleClose}>Close</button>
			<button type="submit">Next</button>
			<div>
				<img src={`${media.media_src.large}`} alt="" />
				<p>{caption}</p>
				<p>{datKey}</p>
				<p>{location}</p>
			</div>
			<div>
				<button type="submit">Add Caption</button>
				<button type="submit">Add Location</button>
				<DeleteMediaButton mediaID={media.id} handleDelete={handleDelete} />
			</div>
		</article>
	);
};

interface DeleteMediaButtonProps {
	mediaID: string;
	handleDelete: () => void;
}
const DeleteMediaButton: React.FunctionComponent<
	React.PropsWithChildren<DeleteMediaButtonProps>
> = (props: DeleteMediaButtonProps) => {
	return <button onClick={props.handleDelete}>Delete</button>;
};

export default MediaDetail;
