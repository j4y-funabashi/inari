'use client';

interface CollectionDetailParams {
    params: {
        id: string
    }
}

export default function CollectionDetailPage({ params }: CollectionDetailParams) {

    return (
        <h1>{params.id}</h1>
    )
}
