
export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Inari" },
    { name: "description", content: "" },
  ];
}

export async function clientLoader(): Promise<Collection[]> {
  const url = "/api/timeline/months?type=hashtag"
  const res = await fetch(url)
  return res.json()
}


export default function CollectionDetailPage({ loaderData }: Route.ComponentProps) {
  return <CollectionList data={loaderData} />;
}
