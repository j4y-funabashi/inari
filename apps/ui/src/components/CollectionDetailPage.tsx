import React from "react";
import { TimelineMonthQuery, timelineMonthResponse } from "../apiClient";
import { useParams } from "react-router-dom";
import CollectionDetail from "./CollectionDetail";

type urlParams = {
	collection_type: string
	collection_id: string;
};
interface MediaTimelineMonthPageProps {
	fetchTimelineMonth: TimelineMonthQuery;
}

const CollectionDetailPage: React.FunctionComponent<
	React.PropsWithChildren<MediaTimelineMonthPageProps>
> = (props: MediaTimelineMonthPageProps) => {
	const [timelineData, setTimelineData] = React.useState<timelineMonthResponse>(
		{
			media: [],
			collection_meta: { title: "", id: "", type: "", media_count: 0 },
		},
	);

	const { collection_id } = useParams<urlParams>();
	console.log(collection_id);

	React.useEffect(() => {
		(async () => {
			const timelineResponse = await props.fetchTimelineMonth(collection_id!);
			setTimelineData(timelineResponse);
		})();
	}, [setTimelineData, collection_id, props]);

	console.log(timelineData);

	return <CollectionDetail mediaTimeline={timelineData} />;
};

export default CollectionDetailPage;
