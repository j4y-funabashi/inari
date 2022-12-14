import React from "react";
import { timelineMonthResponse, media, mockMedia } from "../apiClient";
import { format } from "date-fns";
import MediaDetail from "./MediaDetail";

export interface MediaTimelineMonthProps {
	mediaTimeline: timelineMonthResponse;
}

interface currentMedia {
	isVisible: boolean;
	media: media;
}

const MediaTimelineMonth: React.FunctionComponent<
	React.PropsWithChildren<MediaTimelineMonthProps>
> = (props: MediaTimelineMonthProps) => {
	const { mediaTimeline } = props;

	const [currentMedia, setCurrentMedia] = React.useState<currentMedia>({
		isVisible: false,
		media: mockMedia(new Date(1984, 0, 28, 19, 0, 52)),
	});

	// sort by date
	mediaTimeline.media.sort((a, b) => {
		if (a.date < b.date) {
			return -1;
		}
		return 1;
	});
	// group media to days
	let dayCollections = new Map<string, timelineMonthResponse>();
	mediaTimeline.media.forEach((m) => {
		const dat = new Date(m.date);
		const datKey = format(dat, "yyyy-MM-dd");
		const datTitle = format(dat, "eee, do MMM");
		const dayCollection = dayCollections.get(datKey);
		if (dayCollection) {
			dayCollection.media.push(m);
			dayCollection.collection_meta.media_count = dayCollection.media.length;
			dayCollections.set(datKey, dayCollection);
		} else {
			const c: timelineMonthResponse = {
				media: [m],
				collection_meta: {
					id: datKey,
					title: datTitle,
					type: "timeline_day",
					media_count: 1,
				},
			};
			dayCollections.set(datKey, c);
		}
	});

	// render
	const media = Array.from(dayCollections.values()).map((v) => {
		const thumbs = v.media.map((m) => {
			return (
				<li key={m.id}>
					<MediaThumb media={m} setCurrentMedia={setCurrentMedia} />
				</li>
			);
		});
		return (
			<div key={v.collection_meta.id}>
				<h2>{v.collection_meta.title}</h2>
				{thumbs}
			</div>
		);
	});

	const header = (
		<div>
			<h1>{mediaTimeline.collection_meta.title}</h1>
			<small>{mediaTimeline.collection_meta.media_count}</small>
		</div>
	);

	return (
		<div>
			{currentMedia.isVisible && (
				<MediaDetail
					media={currentMedia.media}
					handleDelete={() => {
						setCurrentMedia({
							isVisible: false,
							media: currentMedia.media,
						});
					}}
				/>
			)}
			{header}
			{media}
		</div>
	);
};

interface MediaThumbProps {
	media: media;
	setCurrentMedia: React.Dispatch<React.SetStateAction<currentMedia>>;
}
const MediaThumb: React.FunctionComponent<
	React.PropsWithChildren<MediaThumbProps>
> = (props: MediaThumbProps) => {
	const { media } = props;
	const handleMediaClick = () => {
		console.log(`HELL! ${media.id}`);
		props.setCurrentMedia({
			media: media,
			isVisible: true,
		});
	};

	return (
		// rome-ignore lint/a11y/useKeyWithClickEvents: <explanation>
		<img
			alt=""
			width="150"
			src={`${media.media_src.small}`}
			onClick={handleMediaClick}
		/>
	);
};

export default MediaTimelineMonth;
