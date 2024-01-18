import { useState } from "react";
import slugify from "slugify";

interface AddHashtagFormProps {
  handleSaveHashtag: (hashtag: string) => Promise<void>;
}

export const AddHashtagForm = ({ handleSaveHashtag }: AddHashtagFormProps) => {
  const [newText, setNewText] = useState("");

  return (
    <div className="grid grid-cols-6">
      <label>
        Add Hashtag
        <input
          type="text"
          value={newText}
          onChange={(e) => {
            const t = slugify(e.target.value, { strict: true, lower: true });
            setNewText(t);
          }}
          className="col-span-5 text-black"
        />
      </label>
      <button
        className="bg-green text-white font-bold py-1 px-2 col-span-1"
        onClick={async () => {
          if (newText === "") {
            return;
          }
          await handleSaveHashtag(newText);
        }}
      >
        Save
      </button>
    </div>
  );
};
