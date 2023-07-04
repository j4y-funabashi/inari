import { render, screen } from '@testing-library/react'
import { MediaCard } from "@/app/components/mediaCard";
import '@testing-library/jest-dom'
import { Media } from '@/app/apiClient';

it('renders a heading', () => {

    const handleDelete = () => { }
    const m: Media = {
        id: "",
        thumbnails: {
            medium: "md_test-image-123.jpg",
            large: "lg_test-image-123.jpg"
        },
        collections: [
            {
                id: "test-collection-123",
                title: "Test collection 123",
                media_count: 10,
                type: "test-collection"
            }
        ],
        date: "2023-01-23T10:00:00Z",
        caption: "test caption"
    }

    // ARRANGE
    render(<MediaCard m={m} handleDelete={handleDelete} />)

    // ASSERT
    const img = screen.getByRole('img')
    expect(img).toHaveAccessibleName('test caption')
})
