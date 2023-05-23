import format from "date-fns/format";
import React from "react";
import { media } from "../apiClient";

interface MediaProps {
	media: media;
	handleDelete: () => void;
}

const Media: React.FunctionComponent<React.PropsWithChildren<MediaProps>> = (
	props: MediaProps,
) => {
	const { media, handleDelete } = props;

	const dat = new Date(media.date);
	const datKey = format(dat, "eee, do MMM yyyy - HH:mm");
	const thumbnailLink = "/thumbnails/" + media.thumbnails.medium
	// const location = `${media.location.locality}, ${media.location.region}, ${media.location.country.long}`;
	const caption = media.caption;
	const collectionLinks = media.collections.map(
		(c) => {
			const url = "/collection/" + c.type + "/" + c.id
			return (
				<a href={url}>{c.title}</a>
			)
		}
	)

	return (
		<article>
			<div>
				<img src={thumbnailLink} alt="" />
			</div>
			<div>
				<p>{caption}</p>
				<p>{datKey}</p>
				{/* <p>{location}</p> */}
			</div>
			<div>{collectionLinks}</div>
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

export default Media;
