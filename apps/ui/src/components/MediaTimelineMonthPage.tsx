import React from "react";
import { TimelineMonthQuery, timelineMonthResponse } from "../apiClient";
import { useParams } from "react-router-dom";
import MediaTimelineMonth from "./MediaTimelineMonth";

type urlParams = {
	monthid: string;
};
interface MediaTimelineMonthPageProps {
	fetchTimelineMonth: TimelineMonthQuery;
}

const MediaTimelineMonthPage: React.FunctionComponent<React.PropsWithChildren<MediaTimelineMonthPageProps>> = (props: MediaTimelineMonthPageProps) => {
	const [timelineData, setTimelineData] = React.useState<timelineMonthResponse>(
		{
			media: [],
			collection_meta: { title: "", id: "", type: "", media_count: 0 },
		},
	);

	const { monthid } = useParams<urlParams>();
	console.log(monthid);

	React.useEffect(() => {
		(async () => {
			const timelineResponse = await props.fetchTimelineMonth(monthid);
			setTimelineData(timelineResponse);
		})();
	}, [setTimelineData, monthid, props]);

	console.log(timelineData);

	return <MediaTimelineMonth mediaTimeline={timelineData} />;
};

export default MediaTimelineMonthPage;
