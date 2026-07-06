import { listCollections, mockListCollections, type Collection } from "~/apiClient";
import type { Route } from "./+types/home";
import { CollectionList } from "~/components/collectionList";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export async function clientLoader(): Promise<Collection[]> {
  if (process.env.NODE_ENV == "development") {
    return mockListCollections("inbox")
  }
  return await listCollections("inbox")
}

export default function Home({ loaderData }: Route.ComponentProps) {
  return <CollectionList data={loaderData} />;
}
